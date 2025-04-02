/*
Copyright © 2025 Quentin ROBCIS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package models

import "qrobcis/pkgsmanager/internal/types/provider"

type RawPackageConfiguration struct {
    Name       string `yaml:"name"`
    GPGKey     string `yaml:"gpgKey"`
    SourceList string `yaml:"souceList"`
    Provider   string `yaml:"provider"`
    Version    string `yaml:"version"`
}

type PackageConfiguration struct {
    Name       string            `yaml:"name"`
    GPGKey     string            `yaml:"gpgKey"`
    SourceList string            `yaml:"souceList"`
    Provider   provider.Provider `yaml:"provider"`
    Version    string            `yaml:"version"`
}

func NewPackageConfiguration(name string, gpgKey string, souceList string, providerName string, version string) *PackageConfiguration {
    providerValue := provider.ToProvider(providerName)
    if providerValue == provider.Unset {
        providerValue = provider.APT
    }

    return &PackageConfiguration{
        GPGKey:     gpgKey,
        Name:       name,
        Provider:   providerValue,
        Version:    version,
        SourceList: souceList,
    }
}
