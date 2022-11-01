#Author: KC Shasteen
#Title: Earth-Mars Transmission Delay program
#Date: 7/18/2021
#Purpose: This program contains functions meant to be used on an email sever for a Mars analog.  
# These functions will allow the server to calculate the Earth-Mars communications delay for any given 
# point in a hypothetical mission context.  This delay can then be added to a message before delevery.
# This will allow the server to more realistically approximate the experience of communicating with
# Earth from Mars while inside a Mars analogue base.
# The program also contains code for displaying useful information about how the time delay changes
# over time and where the inner planets are in their orbits.
# Ephemerides are taken from NASA JPL Horizons Telnet interface for Epoch J2000.
# Code has been ported to Python from an OpenGL model solar system implementation by same author.

import numpy as np
import math 
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
import datetime

#define KMOVERAU 149598000  #there are 149,598,000 km per au

#This subfunction is used to solve Keplers equation of equal areas in equal times.
##inputs are:
#x, an x-coordinate within an ellipse with origin centered at the left focus and with semimajor axis along the x-axis
#a, the semimajor axis length, and
#b, the semiminor axis length, such that a >= b
#returns the area between the x-axis and the top half of the ellipse that is left of a line drawn from the left focus to a point on the ellipse at coordinate x
def ellipseAreaFunc(x, a, b):
    asquared = a * a
    ecc = math.sqrt(1.0 - b * b / (asquared))
    abovertwo = a * b / 2.0

    for i in range (0,x.size,1): #ensures that x is not larger than than samimajor axis
        if (x[i] > a):
            x[i] = a
        if (x[i] < -a):
            x[i] = -a

    return (abovertwo*(np.pi - np.arccos(x / a) - ecc * np.sqrt(1.0 - x * x / (asquared))))

#This function solves Keplers equation of equal areas in equal times in order to return coordintates in x,y space.
#inputs are:
#t, the amount of time that has passed since a periapsis
#p, the amount of time for from one periapsis to the next
#a, the semimajor axis length, and
#b, the semiminor axis length, such that a >= b
#returns the x and y coordinates of an oribiting body within an orbital ellipse with origin centered at the left focus and with semimajor axis along the x-axis
#solves in to adequate precision in about 32 steps using a simple binary search with no convergence accelerations
def getOrbitalCoords(t, period, a, b):
    bottomHalf = False
    x = 0.0
    y = 0.0
    area = ellipseAreaFunc(np.array([x]), a, b)[0]
    areaTarget = np.pi * a*b*(t - math.floor(t / period)*period) / period #equal areas in equal times...
    # floor(t/period) == number of full periods that have passed since t = 0
    # floor(t/period)*period == total time it took for that many periods to pass
    # t - floor(t/period)*period == the time ellapsed since the beginning of most recent period
    # areaTarget == the total area swept out by a line from the sun to the planet since the beginning of the most recent period where the beginning of the period is defined as the point in time where the planet was at perihelion

    if (not(areaTarget < np.pi*a*b / 2.0)): #this could probably be simplified with if (t < period /2):
       areaTarget = np.pi * a*b - areaTarget
       bottomHalf = True
    
    pp = .5
    while (abs(area - areaTarget) > .0000001):
        area = ellipseAreaFunc(np.array([x]), a, b)[0]

        if (area > areaTarget):
            x -= a * pp;
        if (area < areaTarget):
            x += a * pp;

        pp *= .5;

    if (x > a):
        x = a;
    if (x < -a):
        x = -a;

    y = b * math.sqrt(1.0 - x * x / (a*a));

    if (bottomHalf):
        y *= -1.0;    

    return (x, y)

#inputs:
#angle, the inital angle (in degrees) of the orbital body in its orbital ellipse with origin centered at the left focus and with semimajor axis along the x-axis, where 0° is periapsis, and 180 is apoapsis
#a, length of the semimajor axis
#e, eccentricity of the ellipse
#returns a coefficient that when multiplied by the period produces the amount of time that would have passed for the planet to reach its initial position from perihelion
def findInitialTemporalDisplacement(angle, a, e):
    #ensures 0 < angle < 360
    while (angle < 0.0): 
        angle += 360.0
    while (angle > 360.0):
        angle -= 360.0
    #t;
    #x;
    #area;
    b = a * math.sqrt(1.0 - e * e);

    if (angle != 180.0):
        x = a * b*b*(e - (1.0 / math.cos(angle*np.pi / 180.0))) / (b*b + a * a*math.tan(angle*np.pi / 180.0)*math.tan(angle*np.pi / 180.0));

        if (angle < 180.0):
            area = (a*b / 2.0)*(np.pi - math.acos((x - e * a) / a) - e * math.sqrt(1.0 - (x - e * a)*(x - e * a) / (a*a)));
        else:
            area = np.pi * a*b - (a*b / 2.0)*(np.pi - math.acos((x - e * a) / a) - e * math.sqrt(1.0 - (x - e * a)*(x - e * a) / (a*a)));
    else:
        area = np.pi * a*b / 2.0;

    if (area > np.pi*a*b):
        area = np.pi * a*b;
    if (area < 0.0):
        area = 0.0;

    t = area / (np.pi*a*b);

    return (t); #returns a coefficient that when multiplied by the period produces the amount of time that would have passed for the planet to reach its initial position from perihelion

#//calculates the true longitude from the meanlongitdue, the longitude of perihelion, and the orbital eccentricity by iteration
#// true longitude = longitude of periapsis + true anomoly ;
#// mean anomoly = eccentric anomoly - eccentricity * sin (eccentric anomoly) ;
#// true anomoly = arccos ( (cos(eccentric anomoly) - eccentricity)/(1-eccentricity*cos(eccentric anomoly)) ) ;
#// mean longitude = mean anomoly + longitude of periapsis

#// t = true longitude
#// a = true anomoly
#// M = mean anomoly
#// L = mean longitude
#// q = eccentric anomoly
#// p = longitude of periapsis
#// e = eccentricity

#// t = p + a
#// M = q - e * sin(q)
#// a = arccos( (cos(q) - e) / (1 - e * cos(q)) )
#// L = M + p
def getTrueLongitude(mlg, lop, ecc):
    while ((mlg < 0.0)):
        mlg += 360.0;

    while (not(mlg < 360.0)):
        mlg -= 360.0;

    #//find eccentric anomoly by iteration
    #//meanAnom(q) = q - e * sin(q)
    q = np.pi;
    targetMean = (mlg - lop)*np.pi / 180.0;#//mean anomoly in radians

    while ((targetMean < 0.0)):
        targetMean += 2.0*np.pi;

    while (not(targetMean < 2.0*np.pi)):
        targetMean -= 2.0*np.pi;

    mean = q - ecc * math.sin(q); #//guessed mean from guessed q

    pp = .5;

    while (abs(mean - targetMean) > .0000001):
        mean = q - ecc * math.sin(q);

        if (mean > targetMean):
            q -= np.pi * pp;
        if (mean < targetMean):
            q += np.pi * pp;

        pp *= .5;

    #//use, q, the eccentric anomoly to find t, the true longitude
    #// t = p + arccos( (cos(q) - e) / (1 - e * cos(q)) )
    #// t = lop + acos( (cos(q) - ecc) / (1.0 - ecc * cos(q)) ) * 180.0/PI
    trueLong = 0.0;
    if (q > np.pi):
        trueLong = float (360.0 - (-lop + math.acos((math.cos(q) - ecc) / (1.0 - ecc * math.cos(q))) * 180.0 / np.pi));
    else:
        trueLong = float (lop + math.acos((math.cos(q) - ecc) / (1.0 - ecc * math.cos(q))) * 180.0 / np.pi);

    while ((trueLong < 0.0)):
        trueLong += 360.0;

    while (not(trueLong < 360.0)):
        trueLong -= 360.0;

    return (trueLong); #//trueLong;

#transforms the coordinates within an orbital ellipse with origin centered at the left focus and with semimajor axis along the x-axis to a right handed 3d solarsystem coordinates with positive x-axis along the first point of ares for epoch J2000
#inputs:
#vecs, the original coordinates
#loa, longitude of ascention 
#inc, inclination
#lop, longitude of periapsis
#ecc, eccentricity
#sma, semimajor axis length
#returns the transformed coordinates
def transformCoords(vecs,loa,inc,lop,ecc,sma):
    from scipy.spatial.transform import Rotation
    #0th get input coords
    #1st translate to focus
    #2nd flip foci
    #3rd rotate to perihelion
    #4th rotate about the line of nodes

    #vector along the line of nodes
    vec = [math.cos(loa * np.pi / 180.0), math.sin(loa * np.pi / 180.0), 0.0]
    vec /= np.linalg.norm(vec) #normalizes vector
    vec *= inc * np.pi / 180.0 #makes magnitude of vector = rotation angle in radians
    rotNode = Rotation.from_rotvec(vec)

    #glRotatef(lop, 0.0, 0.0, 1.0); //orients orbit to perihelion
    vec = [0.0, 0.0, 1.0]
    vec /= np.linalg.norm(vec) #normalizes vector
    vec *= lop * np.pi / 180.0 #makes magnitude of vector = rotation angle in radians
    rotPeri = Rotation.from_rotvec(vec)

    #glRotatef(180.0, 0.0, 0.0, 1.0); //flips positions of right and left focus
    vec = [0.0, 0.0, 1.0]
    vec /= np.linalg.norm(vec) #normalizes vector
    vec *= 180.0 * np.pi / 180.0 #makes magnitude of vector = rotation angle in radians
    rotFlip = Rotation.from_rotvec(vec)

    #combine previous rotations, starting with rotFlip, ending with rot Node   
    rotAll = rotNode * rotPeri * rotFlip 

    #glTranslatef(ecc*sma, 0.0, 0.0); //places center of orbit on left focus
    #glTranslatef(planetXTrans, -planetYTrans, 0.0); //positions planet in orbit for given input coords
    vec2 = np.array([ecc*sma,0,0])
    newvecs = vecs.transpose() + vec2 #translates the initial vectors

    newvecs = rotAll.apply(newvecs) #applies the rotations
    newT = newvecs.transpose() #returns the transpose

    return (newT) 

#//source: nasa fact sheets (except RGB values, those are eyeballed from wikipedia pictures of each planet)
#vectors are J2000 epoch
earthTilt = 23.45;
mercuryPolarVector = [222.8756674099116,   -199.1603624915450,      2421.321473263892]
venusPolarVector   = [113.1130691873618,    65.79833385298434,      6050.385049965100]
earthPolarVector   = [5.644067347543868,    2528.725504468390,      5832.136439780563]
marsPolarVector    = [1506.319358939092,   -187.4395033575126,      3015.724599738247]

planets = [ ["Mercury", "CELESTIALBODY_CRATERED", 0.8, 0.8, 0.8, 58.6462, 87.969257, 0.38709893, 0.20563069, 77.45645, 7.00487, 48.33167, 252.25084, 0.0, mercuryPolarVector, 2439.7, 2439.7, "mercuryTextureNumber"],
            ["Venus", "CELESTIALBODY_CLOUDY", .965, .965, .573, -243.0185, 224.70079922, 0.72333199, 0.00677323, 131.53298, 3.39471, 76.68069, 181.97973, 177.36, venusPolarVector, 6051.8, 6051.8, "venusTextureNumber"],
            ["Earth", "CELESTIALBODY_ICECAPPED", 0.0, .098, 1.0, 0.997257917, 365.25636, 1.00000011, 0.01671022, 102.94719, 0.00005, -11.26064, 100.46435, earthTilt, earthPolarVector, 6378.1, 6356.8, "earthTextureNumber"],
            ["Mars", "CELESTIALBODY_ICECAPPED", .541, 0.0, 0.0, 1.02595675, 686.98, 1.52366231, 0.09341233, 336.04084, 1.85061, 49.57854, 355.45332, 25.19, marsPolarVector, 3396.2, 3376.2, "marsTextureNumber"] ]

#index variable             units or range
#0  name                    str
#1  appearance              str
#2  red                     0-1.0
#3  grn                     0-1.0
#4  blu                     0-1.0
#5  rotation                (days)
#6  revolution              (days)
#7  semimajoraxis           (AU)
#8  orbitalEccentricty      0-1.0
#9  longitudeOfPerihelion   (deg)
#10 orbitalInclination      (deg)
#11 longitudeOfAcension     (deg)
#12 meanLongitude           (deg)
#13 axialTilt               (deg)
#14 polarVector             {X Y Z}
#15 equatorialRadius        km
#16 polarRadius             km
#17 textureNumber           N/A

#This function calculates the light lag from Earth to Mars.
#inputs:
#date, in datetime datastructure
#showplot, if True displays a small map of innner solarsystem with planetary positions and orbits, if False this is skipped
#verbose, if True prints various useful figures to the terminal, if False this is skipped
#returns Earth-Mars communications delay in seconds
def getEarthMarsDelay(desiredDate=datetime.datetime.now(),showplot=False,verbose=False):
    today = datetime.date.today()
    if verbose:
        print("Today's date:", today)

    epoch = datetime.datetime(2000, 1, 1, 12, 0)        # J2000 epoc, starting date for planets based on ephemerides
    #desiredDate  = datetime.datetime.now()              # Now
    #desiredDate = datetime.datetime(2033, 7, 5, 12, 0) # Uncomment this to enter a specific date (a transit or conjunction for example)
    duration = desiredDate - epoch                      # varible holds a date object with the differnce in time from desired date to epoch
    delta = duration.total_seconds()                    # Total number of seconds between dates

    timeElapsed = delta/(24*60*60) #days since j2000 epoch
    if verbose:
        print ("Showing planetary positions for: ", datetime.datetime(2000, 1, 1, 12, 0) + datetime.timedelta(days=timeElapsed) )

    if showplot:
        fig = plt.figure(figsize=(9.5, 9.5))
        ax = fig.add_subplot(1, 1, (1, 2), projection='3d') #rows, columns, postion(1=top left, increases to right)

        ax.plot(np.array([0,planets[3][7]]),np.array([0,0]),np.zeros(2), label='First Point of Ares')
    planetCoordsAtEllapsedTime = ([[0,0,0],[0,0,0],[0,0,0],[0,0,0]])

    for current in range (0,4,1):
        mlg = planets[current][12]
        loa = planets[current][11]
        inc = planets[current][10]
        lop = planets[current][9]
        ecc = planets[current][8]
        sma = planets[current][7]
        smb = (sma * math.sqrt(1.0 - ecc * ecc)) #//e = sqrt(1-b*b/(a*a))  ->  b*b= a*a - a*a*e*e  ->  b = a * sqrt( 1 - e*e )
        rev = planets[current][6]
        itd = findInitialTemporalDisplacement(getTrueLongitude(mlg, lop, ecc) - lop, sma, ecc)
        itp = initialTemporalPhase = rev * itd
        planetXTrans, planetYTrans = getOrbitalCoords(-timeElapsed - itp, rev, sma, smb)

        if showplot:
            #calcultes orbit coordinates
            size = 100; xcoord = np.zeros(size+1); ycoord = np.zeros(size+1)
            for i in range (0,size+1):
                coords = getOrbitalCoords(-i/size,1,sma,smb)
                xcoord[i] = coords[0]; ycoord[i] = coords[1]
            zcoord=np.zeros(size+1)

            #places and draws the orbits
            vecs = np.array([xcoord,ycoord,zcoord])
            newT = transformCoords(vecs,loa,inc,lop,ecc,sma)
            ax.plot(newT[0],newT[1],newT[2], label=planets[current][0],color=(planets[current][2],planets[current][3],planets[current][4]),linewidth=3)#, marker='x',markersize=16)

        #places and draws the planets
        vecs = np.array([planetXTrans,planetYTrans,0.0])
        newT = transformCoords(vecs,loa,inc,lop,ecc,sma)
        if showplot:
            ax.plot(np.array([newT[0]]),np.array([newT[1]]),np.array([newT[2]]), marker='o',markersize=16,color=(planets[current][2],planets[current][3],planets[current][4]))

        #stores planet positions for later use
        planetCoordsAtEllapsedTime[current] = newT

    pcet = planetCoordsAtEllapsedTime #too long
    if verbose:
        for i in range (0,4):   
            print (planets[i][0],pcet[i])

    #calculates distance and delay
    EMdist = math.sqrt((pcet[2][0]-pcet[3][0])**2+(pcet[2][1]-pcet[3][1])**2+(pcet[2][2]-pcet[3][2])**2) #in AU
    EMdistKM = EMdist*149598000 #in km; 149598000 km / au
    EMdelay = EMdist*149598000000/299792458 #in seconds; light speed: 299 792 458 m / s
    if verbose:
        print('Earth-Mars time delay:',EMdelay,'sec')

    if showplot:
        #displays results
        fig.text(.05, .925, s='Planets Shown for Date: %s'%desiredDate, fontsize=12,color='white')
        fig.text(.05, .90, s='Distance from Earth to Mars: %f [AU]'%EMdist, fontsize=12,color='white')
        fig.text(.05, .875, s='Distance from Earth to Mars: %f [km]'%EMdistKM, fontsize=12,color='white') 
        fig.text(.05, .85, s='Time Delay from Earth to Mars: %i [min] %.1f [s]'%(math.floor(EMdelay/60),(EMdelay%60)), fontsize=12,color='white') 
        fig.text(.05, .825, s='Two-way Delay (Earth to Mars): %i [min] %.1f [s]'%(math.floor(2*EMdelay/60),(2*EMdelay%60)), fontsize=12,color='white') 

        #configures background, view, and axes
        fig.patch.set_facecolor('black')
        ax.set_xlabel('x [AU]'); ax.set_xlim([-sma*1.1, sma*1.1])
        ax.set_ylabel('y [AU]'); ax.set_ylim([-sma*1.1, sma*1.1])
        ax.set_zlabel('z [AU]'); ax.set_zlim([-sma*1.1, sma*1.1])
        ax.view_init(elev=90, azim=-90)
        ax.set_facecolor('black')
        ax.xaxis.label.set_color('red')
        ax.yaxis.label.set_color('red')
        ax.zaxis.label.set_color('red')
        ax.tick_params(axis='x', colors='red')
        ax.tick_params(axis='y', colors='red')
        ax.tick_params(axis='z', colors='red')
        ax.legend(framealpha=0.5)

        plt.show()
    return EMdelay #earth mars delay in seconds

############### Show Orbits and Time Delay ###################
# This segment displays a graph of the planets, the date of interest, and the 1 way and 2 way Earth Mars communications delay.
# The function will also return the time delay in sections for use in an analogue email server.
desiredDate = datetime.datetime.now() 
#desiredDate = datetime.datetime(2033, 5, 1, 12, 0)  #use this to select a specific date for display
getEarthMarsDelay(desiredDate,showplot=True) 
################# End Show Orbits and Time Delay #############

############### Show Mission Profile ###################
#This segment displays a graph showing an ideal launch window in 2033 along with a hypothetical mission profile with trajectories from:
# Wooster et al. Trajectory Options for Human Mars Missions. American Institute of Aeronautics and Astronautics. 2006. ( https://smartech.gatech.edu/handle/1853/14747 )
startDate = datetime.datetime.now()  #datetime.datetime(2033, 5, 1, 12, 0)
xarray = []
yarray = []
for i in range (0,34):
    desiredDate = startDate + i*datetime.timedelta(days=30)
    xarray.append(desiredDate.strftime('%m/%d/%Y'))
    delay = getEarthMarsDelay(desiredDate,showplot=False)
    yarray.append(delay)
    #print (desiredDate,'|',delay)
    
print (xarray)
print (yarray)

fig = plt.figure(figsize=(9.5, 9.5))
ax = fig.add_subplot(1, 1, (1, 2)) #rows, columns, postion(1=top left, increases to right)

fig.text(0.1, .92, 'Earth-Mars Signal Delay During Year 2033 Conjunction', fontsize=14)
fig.text(0.1, .89, 'Mission profile uses 2-year free return trajectory from Wooster et al. (2006)', fontsize=12)
ax.plot(xarray,np.array(yarray)/60,label='Earth-Mars delay [mins]')
ax.set_xticklabels(xarray, rotation =90,fontsize=8)
ax.axvline(x=0.0, alpha = 0.5,color='#888888') #earth departure (2 year free return trajectory)
ax.axvline(x=3.5, alpha = 0.5,color='#888888') #mars arrival (t=105 days)
ax.axvline(x=24.63, alpha = 0.5,color='#888888') #mars departure after 634 day stay (t = 739 days)
ax.axvline(x=33.6, alpha = 0.5,color='#888888') #earth return 270 day return trip (t=1008 days)
ax.text(0.00+.2, 8.5, 'Earth departure (2 year free return trajectory)', fontsize=10,rotation =90,color='#888888')
ax.text(3.5+.2, 10.5, 'Mars arrival (t=105 days)', fontsize=10,rotation =90,color='#888888')
ax.text(24.63+.2, 8.5, 'Mars departure after 634 day stay (t = 739 days)', fontsize=10,rotation =90,color='#888888')
ax.text(33.6+.2, 8.5, 'Earth return 270 day return trip (t=1008 days)', fontsize=10,rotation =90,color='#888888')
ax.set_xlabel('Date');
ax.set_ylabel('Earth-Mars delay [mins]');
ax.legend(framealpha=0.5)
plt.show()
################# End Mission Profile ##################