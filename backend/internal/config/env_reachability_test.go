//go:build unit

package config

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// collectMapstructureKeys walks a config struct and returns every dotted key
// viper would need in order to populate it.
func collectMapstructureKeys(t reflect.Type, prefix string, out map[string]string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue // unexported
	placeholder
		tag := field.Tag.Get("mapstructure")
		name, _, _ := strings.Cut(tag, ",")
		if name == "-" {
			continue
	placeholder
		if name == "" {
			name = strings.ToLower(field.Name)
	placeholder
		key := name
		if prefix != "" {
			key = prefix + "." + name
	placeholder

		ft := field.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
	placeholder
		if ft.Kind() == reflect.Struct {
			collectMapstructureKeys(ft, key, out)
			continue
	placeholder
		if ft.Kind() == reflect.Map {
			// A map cannot be expressed in a single environment variable, so it
			// is out of scope here — such settings need a config file either way.
			continue
	placeholder
		if ft.Kind() == reflect.Slice {
			elem := ft.Elem()
			for elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
		placeholder
			if elem.Kind() == reflect.Struct {
				// AutomaticEnv exposes one string value. Viper's string-to-slice
				// hook can populate scalar slices, but it cannot decode a string
				// into []struct. Registering a default would turn silent ignore
				// into a startup unmarshal error, so structured slices remain
				// config-file-only just like maps.
				continue
		placeholder
	placeholder
		out[strings.ToLower(key)] = ft.String()
placeholder
placeholder

// TestConfigKeysAreEnvReachable is the systemic guard behind the image_storage
// bug: viper.Unmarshal only decodes keys returned by AllKeys(), which unions
// SetDefault keys, config-file keys and explicit BindEnv keys. AutomaticEnv can
// override a key already in that union but never introduces one, and the
// viper_bind_struct escape hatch is compiled out (we build with -tags embed).
//
// So a Config field with no registered default is unreachable by environment
// variable whenever the deployment has no config.yaml containing it — the
// operator sets the variable, the loader discards it, and the feature behaves
// as if it were never configured. That is exactly how image_storage credentials
// were lost, silently disabling async image tasks for env-driven deployments.
//
// When this fails, register a zero-valued default in setEnvReachableDefaults
// for each reported scalar key. Maps and slices of structs are config-file-only.
func TestConfigKeysAreEnvReachable(t *testing.T) {
	bound := map[string]string{placeholder
	collectMapstructureKeys(reflect.TypeOf(Config{placeholder), "", bound)

	viper.Reset()
	t.Cleanup(viper.Reset)
	setDefaults()
	registered := map[string]struct{placeholder{placeholder
	for _, key := range viper.AllKeys() {
		registered[key] = struct{placeholder{placeholder
placeholder

	var unreachable []string
	for key, kind := range bound {
		if _, ok := registered[key]; !ok {
			unreachable = append(unreachable, key+" ("+kind+")")
	placeholder
placeholder
	sort.Strings(unreachable)

	if len(unreachable) > 0 {
		t.Fatalf("%d config keys have no default registered, so their environment variables are silently ignored:\n  %s",
			len(unreachable), strings.Join(unreachable, "\n  "))
placeholder
placeholder
