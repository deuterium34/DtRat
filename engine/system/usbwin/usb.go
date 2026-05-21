// Gemeni

// Package usb-win предоставляет инструменты для управления USB-устройствами в Windows
// исключительно средствами Pure Go (без CGO) через SetupAPI.
package usbwin

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"unsafe"
)

// Определение структур и констант WinAPI
const (
	DIGCF_PRESENT      = 0x00000002
	DIGCF_ALLCLASSES   = 0x00000004
	SPDRP_DEVICEDESC   = 0x00000000
	SPDRP_FRIENDLYNAME = 0x0000000C

	DIF_PROPERTYCHANGE       = 0x00000012
	DICS_ENABLE              = 0x00000001
	DICS_DISABLE             = 0x00000002
	DICS_FLAG_GLOBAL         = 0x00000001
	DICS_FLAG_CONFIGSPECIFIC = 0x00000002

	DN_HAS_PROBLEM   = 0x00000400
	CM_PROB_DISABLED = 22 // Код проблемы: устройство отключено

	ERROR_NO_MORE_ITEMS = 259
)

// GUID структура для WinAPI
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

// SP_DEVINFO_DATA содержит информацию об устройстве.
// Размер структуры вычисляется динамически для совместимости 32/64 bit
type SP_DEVINFO_DATA struct {
	cbSize    uint32
	ClassGuid GUID
	DevInst   uint32
	Reserved  uintptr
}

// SP_CLASSINSTALL_HEADER — заголовок для параметров установки
type SP_CLASSINSTALL_HEADER struct {
	cbSize          uint32
	InstallFunction uint32
}

// SP_PROPCHANGE_PARAMS — параметры для изменения состояния устройства
type SP_PROPCHANGE_PARAMS struct {
	ClassInstallHeader SP_CLASSINSTALL_HEADER
	StateChange        uint32
	Scope              uint32
	HwProfile          uint32
}

// Загрузка DLL и функций WinAPI
var (
	modsetupapi = syscall.NewLazyDLL("setupapi.dll")
	modcfgmgr32 = syscall.NewLazyDLL("cfgmgr32.dll")

	procSetupDiGetClassDevs              = modsetupapi.NewProc("SetupDiGetClassDevsW")
	procSetupDiEnumDeviceInfo            = modsetupapi.NewProc("SetupDiEnumDeviceInfo")
	procSetupDiGetDeviceRegistryProperty = modsetupapi.NewProc("SetupDiGetDeviceRegistryPropertyW")
	procSetupDiGetDeviceInstanceId       = modsetupapi.NewProc("SetupDiGetDeviceInstanceIdW")
	procSetupDiDestroyDeviceInfoList     = modsetupapi.NewProc("SetupDiDestroyDeviceInfoList")
	procSetupDiSetClassInstallParams     = modsetupapi.NewProc("SetupDiSetClassInstallParamsW")
	procSetupDiCallClassInstaller        = modsetupapi.NewProc("SetupDiCallClassInstaller")

	procCMGetDevNodeStatus = modcfgmgr32.NewProc("CM_Get_DevNode_Status")
)

// USBDevice описывает подключенное USB-устройство
type USBDevice struct {
	Name      string
	DeviceID  string
	VendorID  string
	ProductID string
	IsActive  bool
}

// GetUSBDevices возвращает список всех USB-устройств в системе (активных и отключенных)
func GetUSBDevices() ([]USBDevice, error) {
	var devices []USBDevice

	// Указываем энумератор "USB" для фильтрации
	enumerator, err := syscall.UTF16PtrFromString("USB")
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации строки: %v", err)
	}

	// Получаем дескриптор списка устройств
	hDevInfo, _, err := procSetupDiGetClassDevs.Call(
		0,
		uintptr(unsafe.Pointer(enumerator)),
		0,
		DIGCF_ALLCLASSES, // Берем все устройства (даже отключенные)
	)
	// Syscall возвращает handle или INVALID_HANDLE_VALUE (-1)
	if hDevInfo == ^uintptr(0) {
		return nil, fmt.Errorf("SetupDiGetClassDevsW ошибка: %v", err)
	}
	defer procSetupDiDestroyDeviceInfoList.Call(hDevInfo)

	devInfoData := SP_DEVINFO_DATA{
		cbSize: uint32(unsafe.Sizeof(SP_DEVINFO_DATA{})),
	}

	reVidPid := regexp.MustCompile(`(?i)VID_([0-9A-F]{4})&PID_([0-9A-F]{4})`)

	for i := 0; ; i++ {
		r1, _, err := procSetupDiEnumDeviceInfo.Call(
			hDevInfo,
			uintptr(i),
			uintptr(unsafe.Pointer(&devInfoData)),
		)
		if r1 == 0 {
			if err == syscall.Errno(ERROR_NO_MORE_ITEMS) {
				break // Устройства закончились, штатный выход
			}
			return nil, fmt.Errorf("SetupDiEnumDeviceInfo ошибка на индексе %d: %v", i, err)
		}

		// Получаем Instance ID (путь устройства)
		instIdBuf := make([]uint16, 512)
		var reqSize uint32
		procSetupDiGetDeviceInstanceId.Call(
			hDevInfo,
			uintptr(unsafe.Pointer(&devInfoData)),
			uintptr(unsafe.Pointer(&instIdBuf[0])),
			uintptr(len(instIdBuf)),
			uintptr(unsafe.Pointer(&reqSize)),
		)
		deviceID := syscall.UTF16ToString(instIdBuf)

		// Получаем имя (FriendlyName или DeviceDesc)
		name := getDevicePropertyString(hDevInfo, &devInfoData, SPDRP_FRIENDLYNAME)
		if name == "" {
			name = getDevicePropertyString(hDevInfo, &devInfoData, SPDRP_DEVICEDESC)
		}

		// Парсим VID и PID
		vid, pid := "", ""
		matches := reVidPid.FindStringSubmatch(deviceID)
		if len(matches) == 3 {
			vid = matches[1]
			pid = matches[2]
		}

		// Проверяем статус через Cfgmgr32
		isActive := true
		var status, problem uint32
		r1, _, _ = procCMGetDevNodeStatus.Call(
			uintptr(unsafe.Pointer(&status)),
			uintptr(unsafe.Pointer(&problem)),
			uintptr(devInfoData.DevInst),
			0,
		)
		if r1 == 0 { // CR_SUCCESS
			// Если флаг проблемы установлен и это "DISABLED"
			if (status&DN_HAS_PROBLEM) != 0 && problem == CM_PROB_DISABLED {
				isActive = false
			}
		}

		devices = append(devices, USBDevice{
			Name:      name,
			DeviceID:  deviceID,
			VendorID:  vid,
			ProductID: pid,
			IsActive:  isActive,
		})
	}

	return devices, nil
}

// EnableDevice включает устройство по его DeviceInstanceID
func EnableDevice(deviceID string) error {
	return changeDeviceState(deviceID, DICS_ENABLE)
}

// DisableDevice отключает устройство по его DeviceInstanceID
func DisableDevice(deviceID string) error {
	return changeDeviceState(deviceID, DICS_DISABLE)
}

// EnableDevicePnPUtil резервный надежный метод включения через консольную утилиту Windows
func EnableDevicePnPUtil(deviceID string) error {
	cmd := exec.Command("pnputil", "/enable-device", deviceID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pnputil ошибка: %v, вывод: %s", err, string(out))
	}
	return nil
}

// DisableDevicePnPUtil резервный надежный метод отключения через консольную утилиту Windows
func DisableDevicePnPUtil(deviceID string) error {
	cmd := exec.Command("pnputil", "/disable-device", deviceID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pnputil ошибка: %v, вывод: %s", err, string(out))
	}
	return nil
}

// Вспомогательная функция для получения строковых свойств из реестра устройств
func getDevicePropertyString(hDevInfo uintptr, devInfoData *SP_DEVINFO_DATA, property uint32) string {
	buf := make([]uint16, 512)
	var propType, reqSize uint32

	r1, _, _ := procSetupDiGetDeviceRegistryProperty.Call(
		hDevInfo,
		uintptr(unsafe.Pointer(devInfoData)),
		uintptr(property),
		uintptr(unsafe.Pointer(&propType)),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)*2), // размер в байтах
		uintptr(unsafe.Pointer(&reqSize)),
	)
	if r1 != 0 {
		return syscall.UTF16ToString(buf)
	}
	return ""
}

// changeDeviceState осуществляет смену состояния через WinAPI SetupAPI (DICS_ENABLE / DICS_DISABLE)
// ВАЖНО: Требует прав администратора!
func changeDeviceState(deviceID string, state uint32) error {
	enumerator, _ := syscall.UTF16PtrFromString("USB")
	hDevInfo, _, err := procSetupDiGetClassDevs.Call(
		0,
		uintptr(unsafe.Pointer(enumerator)),
		0,
		DIGCF_ALLCLASSES,
	)
	if hDevInfo == ^uintptr(0) {
		return fmt.Errorf("SetupDiGetClassDevsW ошибка: %v", err)
	}
	defer procSetupDiDestroyDeviceInfoList.Call(hDevInfo)

	devInfoData := SP_DEVINFO_DATA{
		cbSize: uint32(unsafe.Sizeof(SP_DEVINFO_DATA{})),
	}

	for i := 0; ; i++ {
		r1, _, err := procSetupDiEnumDeviceInfo.Call(
			hDevInfo,
			uintptr(i),
			uintptr(unsafe.Pointer(&devInfoData)),
		)
		if r1 == 0 {
			if err == syscall.Errno(ERROR_NO_MORE_ITEMS) {
				break
			}
			return err
		}

		instIdBuf := make([]uint16, 512)
		var reqSize uint32
		procSetupDiGetDeviceInstanceId.Call(
			hDevInfo,
			uintptr(unsafe.Pointer(&devInfoData)),
			uintptr(unsafe.Pointer(&instIdBuf[0])),
			uintptr(len(instIdBuf)),
			uintptr(unsafe.Pointer(&reqSize)),
		)
		currentID := syscall.UTF16ToString(instIdBuf)

		// Ищем совпадение по ID
		if strings.EqualFold(currentID, deviceID) {
			// Настраиваем параметры изменения
			propChange := SP_PROPCHANGE_PARAMS{
				ClassInstallHeader: SP_CLASSINSTALL_HEADER{
					cbSize:          uint32(unsafe.Sizeof(SP_CLASSINSTALL_HEADER{})),
					InstallFunction: DIF_PROPERTYCHANGE,
				},
				StateChange: state,
				Scope:       DICS_FLAG_GLOBAL,
				HwProfile:   0,
			}

			// Устанавливаем параметры
			r1, _, err = procSetupDiSetClassInstallParams.Call(
				hDevInfo,
				uintptr(unsafe.Pointer(&devInfoData)),
				uintptr(unsafe.Pointer(&propChange)),
				uintptr(unsafe.Sizeof(propChange)),
			)
			if r1 == 0 {
				return fmt.Errorf("SetupDiSetClassInstallParams ошибка (возможно, нет прав админа): %v", err)
			}

			// Применяем изменения (вызываем инсталлер)
			r1, _, err = procSetupDiCallClassInstaller.Call(
				DIF_PROPERTYCHANGE,
				hDevInfo,
				uintptr(unsafe.Pointer(&devInfoData)),
			)
			if r1 == 0 {
				return fmt.Errorf("SetupDiCallClassInstaller ошибка: %v", err)
			}

			return nil // Успех
		}
	}

	return fmt.Errorf("устройство с DeviceID '%s' не найдено", deviceID)
}
