/*
 * Copyright (C) 29/12/2020, 15:42 igors
 * This file is part of helmcertifier.
 *
 * helmcertifier is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * helmcertifier is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with helmcertifier.  If not, see <http://www.gnu.org/licenses/>.
 */

package helmcertifier

type CheckFunc func(uri string) (CheckResult, error)

type Registry map[string]CheckFunc

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) GetCheck(name string) (CheckFunc, bool) {
	v, ok := (*r)[name]
	return v, ok
}

func (r *Registry) AddCheck(name string, checkFunc CheckFunc) {
	(*r)[name] = checkFunc
}
