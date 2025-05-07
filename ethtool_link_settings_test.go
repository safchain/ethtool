package ethtool

import (
	"errors"
	"net"
	"reflect"
	"sort"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

// TestGetLinkSettings attempts to get link settings for all non-loopback interfaces.
// This serves as a basic integration test to ensure the function handles
// real ioctl calls, including potential fallbacks, without crashing.
// Asserting specific values is difficult as capabilities vary widely between interfaces and drivers.
func TestGetLinkSettings(t *testing.T) {
	intfs, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}
	for _, intf := range intfs {
		intfName := intf.Name
		if intfName == "lo" {
			continue
		}
		e, err := NewEthtool() // Assumes NewEthtool() correctly gets a socket fd
		if err != nil {
			// Log error for this interface but continue testing others
			t.Errorf("Failed to create Ethtool instance for interface %s: %v", intfName, err)
			continue
		}
		defer e.Close()

		settings, err := e.GetLinkSettings(intfName)

		// We expect either success (returning some data, possibly indicating 'unknown')
		// or a specific error like EOPNOTSUPP if the interface/driver
		// doesn't support either GLINKSETTINGS or GSET.
		if err != nil {
			var errno syscall.Errno
			if errors.As(err, &errno) {
				if errors.Is(errno, unix.EOPNOTSUPP) {
					// This is an expected outcome for some interfaces/drivers
					t.Logf("GetLinkSettings for '%s' returned EOPNOTSUPP (expected for some devices)", intfName)
				} else {
					// Report other unexpected errors
					t.Errorf("GetLinkSettings for '%s' failed with unexpected error: %v", intfName, err)
				}
			} else {
				// Non-syscall error
				t.Errorf("GetLinkSettings for '%s' failed with non-syscall error: %v", intfName, err)
			}
		} else {
			// If successful, check that the settings struct is not nil and Source is set
			if settings == nil {
				t.Errorf("GetLinkSettings for '%s' succeeded but returned nil settings", intfName)
				// Cannot continue checks if settings is nil
				continue
			}
			if settings.Source != SourceGLinkSettings && settings.Source != SourceGSet {
				t.Errorf("GetLinkSettings for '%s' succeeded but Source ('%s') is invalid", intfName, settings.Source)
			}
			t.Logf("GetLinkSettings for '%s' succeeded (Source: %s). Settings: %+v", intfName, settings.Source, settings)

			// Check SupportedLinkModes
			if settings.SupportedLinkModes == nil {
				t.Errorf("GetLinkSettings for '%s' succeeded but SupportedLinkModes is nil", intfName)
			} else {
				t.Logf("Interface '%s' supported modes: %v", intfName, settings.SupportedLinkModes)
			}

			// If source was GSET, verify the conversion logic
			if settings.Source == SourceGSet {
				var cmd EthtoolCmd
				_, errGet := e.CmdGet(&cmd, intfName)
				if errGet != nil {
					t.Errorf("Failed to re-run CmdGet for '%s' to verify GSET conversion: %v", intfName, errGet)
				} else {
					expectedModes := parseLegacyLinkModeMask(cmd.Supported)
					// Sort both slices for reliable comparison
					sort.Strings(expectedModes)
					actualModes := make([]string, len(settings.SupportedLinkModes))
					copy(actualModes, settings.SupportedLinkModes)
					sort.Strings(actualModes)

					if !reflect.DeepEqual(actualModes, expectedModes) {
						t.Errorf("SupportedLinkModes mismatch for '%s' (Source: GSET). Got %v, expected %v (from raw CmdGet)", intfName, actualModes, expectedModes)
					}
				}
			}

			// Basic sanity checks if it succeeded (values might be 0 or unknown)
			if settings.Speed == SPEED_UNKNOWN { // SPEED_UNKNOWN defined in ethtool.go
				t.Logf("Interface '%s' reported speed as SPEED_UNKNOWN (0x%x)", intfName, SPEED_UNKNOWN)
			}
		}
	}
}
