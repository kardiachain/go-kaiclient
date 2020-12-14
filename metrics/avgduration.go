/*
 *  Copyright 2018 KardiaChain
 *  This file is part of the go-kardia library.
 *
 *  The go-kardia library is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Lesser General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  The go-kardia library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  GNU Lesser General Public License for more details.
 *
 *  You should have received a copy of the GNU Lesser General Public License
 *  along with the go-kardia library. If not, see <http://www.gnu.org/licenses/>.
 */

package metrics

import (
	"fmt"
	"time"
)

type AverageDuration struct {
	RollingAvg   int64
	Measurements int64
}

func (p *AverageDuration) Add(duration time.Duration) {
	p.Measurements++

	prevAvg := p.RollingAvg
	p.RollingAvg = prevAvg + (duration.Nanoseconds()-prevAvg)/p.Measurements
}

func (p *AverageDuration) String() string {
	d, _ := time.ParseDuration(fmt.Sprintf("%dns", p.RollingAvg))
	return d.Round(time.Millisecond).String()
}

func (p *AverageDuration) Raw() int64 {
	d, _ := time.ParseDuration(fmt.Sprintf("%dns", p.RollingAvg))
	return d.Round(time.Millisecond).Milliseconds()
}

func (p *AverageDuration) Reset() {
	p.RollingAvg = 0
	p.Measurements = 0
}
