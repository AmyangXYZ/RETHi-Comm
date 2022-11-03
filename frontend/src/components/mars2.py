#!/usr/bin/env python
""" Earth-Mars """
"""
    earth-mars.py - a program to create the Earth Mars distance calculation plus occultation
    Martin J Levy - W6LHI/G8LHI - https://github.com/mahtin/mars-earth
    Copyright (C) 2021 @mahtin - https://github.com/mahtin/earth-mars/blob/main/LICENSE
"""

import math
import time
import datetime
import json
import ephem

class EarthMars(object):
    """ caculate Earth/Mars relationship based on Earth and Mars lat/long """

    def __init__(self, debug=False):
        """ Everything is based on a constellation name """

        self._debug = debug
        self._data = None
        self._earthobserver = None
        self._marsobserver = None

    def __call__(self):
        """ return whatever we have! """
        return self._data

    def earthobserver(self, lon=0.0, lat=0.0, amsl=0.0):
        """ set Earth observer location """

        if lon > 180 or lon < -180 or lat > 90 or lat < -90:
            raise ValueError

        if not self._earthobserver:
            self._earthobserver = ephem.Observer()
        self._earthobserver.lon = self._degrees_to_radians(lon)
        self._earthobserver.lat = self._degrees_to_radians(lat)
        self._earthobserver.elevation = amsl
        return True

    def sun(self, when=None):

        if not self._earthobserver:
            raise RuntimeError

        if when:
            self._earthobserver.date = when
        else:
            # All in UTC - https://rhodesmill.org/pyephem/date.html-
            self._earthobserver.date = datetime.datetime.utcnow()

        sun = ephem.Sun()
        sun.compute(self._earthobserver)

        r = {
            'name': sun.name,
            'mag': sun.mag,
            'radius': self._radians_to_degrees(sun.radius),
            'ra_dec': [self._radians_to_degrees(sun.ra), self._radians_to_degrees(sun.dec)],
            'alt_az': [self._radians_to_degrees(sun.alt), self._radians_to_degrees(sun.az)],
            'earth_distance_km': self._au_to_km(sun.earth_distance),
        }
        return r

    def mars(self, when=None):

        if not self._earthobserver:
            raise RuntimeError

        if when:
            self._earthobserver.date = when
        else:
            # All in UTC - https://rhodesmill.org/pyephem/date.html-
            self._earthobserver.date = datetime.datetime.utcnow()

        mars = ephem.Mars()
        mars.compute(self._earthobserver)

        previous_rising = self._earthobserver.previous_rising(mars).datetime()
        previous_setting = self._earthobserver.previous_setting(mars).datetime()
        next_rising = self._earthobserver.next_rising(mars).datetime()
        next_setting = self._earthobserver.next_setting(mars).datetime()

        r = {
            'name': mars.name,
            'mag': mars.mag,
            'radius': self._radians_to_degrees(mars.radius),
            'ra_dec': [self._radians_to_degrees(mars.ra), self._radians_to_degrees(mars.dec)],
            'alt_az': [self._radians_to_degrees(mars.alt), self._radians_to_degrees(mars.az)],
            'earth_distance_km': self._au_to_km(mars.earth_distance),
            'previous_rising': previous_rising,
            'previous_setting': previous_setting,
            'next_rising': next_rising,
            'next_setting': next_setting,
        }
        return r

    def _au_to_km(self, a):
        """ I think in Km's - even if computers think in AU's """
        return a * 149598073.0

    def _hours_to_degrees(self, h):
        """ I think in degress - even if computers think in radians """
        return h * (360/15)

    def _radians_to_degrees(self, d):
        """ I think in degress - even if computers think in radians """
        return d * (180/math.pi)

    def _degrees_to_radians(self, a):
        """ I think in degress - even if computers think in radians """
        return a / (180/math.pi)

# https://www.techbeamers.com/python-float-range/
import decimal
def float_range(start, stop, step):
    while start < stop:
        yield float(start)
        start += decimal.Decimal(step)

def main():

    m = EarthMars()
    # Santa Cruz (36.9812°N, 122.0262°W)
    m.earthobserver(lon=-122.0262, lat=36.9812, amsl=40.0)
    mars = m.mars()

    now = datetime.datetime.utcnow()

    if False:
        all_times = {}
        all_times['previous_rising'] = mars['previous_rising'].replace(microsecond=0)
        all_times['previous_setting'] = mars['previous_setting'].replace(microsecond=0)
        all_times['now'] = now.replace(microsecond=0)
        all_times['next_rising'] = mars['next_rising'].replace(microsecond=0)
        all_times['next_setting'] = mars['next_setting'].replace(microsecond=0)

        for r in sorted(all_times, key=lambda y: all_times[y]):
            print("%s UTC ; %s" % (all_times[r], r))

    now = now.replace(hour=0, minute=0, second=0, microsecond=0)
    dd = []
    for timeskew in float_range(-27*6, 27*4*5+1, 0.5):
        delta_when = now + datetime.timedelta(weeks=timeskew)
        mars = m.mars(when=delta_when)

        mars_km = mars['earth_distance_km']
        secs = mars_km / (ephem.c / 1000.0)
        delay = datetime.datetime.fromtimestamp(secs)    ## XXX maybe should not use datetime as a way of getting min/sec printed
        dd.append({"date":delta_when.strftime("%Y-%m-%d"), "distance":mars_km, "delay":int(delay.timestamp())})
        sun = m.sun(when=delta_when)
        sun_km = sun['earth_distance_km']

        occult_degrees = 4.0 * (sun['radius'] * 2)

        if abs(mars['ra_dec'][0] - sun['ra_dec'][0]) <= occult_degrees and abs(mars['ra_dec'][1] - sun['ra_dec'][1]) <= occult_degrees:
            if mars_km > sun_km:
                sun_mars = "M"
            else:
                sun_mars = "S"
        else:
            sun_mars = "-"

        # print("%s %12.1f %5s | [%8.2f %8.2f] | [%8.2f %8.2f] | %s | %8.4f < %8.4f ; %6.1f" % (
        #         delta_when, mars_km, delay.strftime("%M:%S"),
        #         mars['ra_dec'][0], mars['ra_dec'][1],
        #         sun['ra_dec'][0], sun['ra_dec'][1],
        #         sun_mars,
        #         mars['radius'],
        #         sun['radius'],
        #         timeskew
        #         )
        # )
    with open("mars-distance.json", 'w') as outfile:
        json.dump(dd, outfile)

if __name__ == '__main__':
    main()