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

import (
    "qrobcis/pkgsmanager/internal/types/provider"
)

type GroupConfiguration struct {
    Name     string
    Packages map[string]*PackageConfiguration
}

func NewGroupConfiguration(name string) *GroupConfiguration {
    return &GroupConfiguration{
        Name:     name,
        Packages: make(map[string]*PackageConfiguration),
    }
}

func (group *GroupConfiguration) AddPackage(configuration *PackageConfiguration) {
    group.Packages[configuration.Name] = configuration
}

func (group *GroupConfiguration) HasProvider(provider provider.Provider) (hasProvider bool) {
    hasProvider = false

    for _, configuration := range group.Packages {
        if configuration.Provider == provider {
            hasProvider = true
            return
        }
    }

    return
}
