

**cropdata.txt**

For following crops the critical S content is given with following references:

SM  silage Maize  (critical S content by Carciotti et al. 2019)
M   maize  (critical S content by Carciotti et al. 2019)
WRA winter rapeseed / Winterraps (critical S content by Ferreira & Ernst 2014)
WW  winter weed (critical S content by Reussi et al. 2012)
SOY soybean (critical S content by Divito et al. 2016)

**All other crops in the list have dummy values. If you want to use them, find the critical S content in the literature and replace the dummy values.**

**irr_sul.txt**

the irrigation file got an extra column for the sulfur content in the irrigation water. The sulfur content is given in mg/l. If you don't have this information, you can set the sulfur content to 0.

**smin.txt**

this file contains measured data of the soil mineral sulfur content. The sulfur content is given in SO4-S kg per ha. 
The initial content is given without date(nan) and the following values are given with the date of the measurement.

**soil_sul.txt**

the soil file contains the c/s ratio (carbon/sulfur ratio) in the organic part of the soil. Normal c/s ratio should be around 100:1.
Please note that the usual n/s ratio is around 8:1, there may be more extreme values in the literature, from 5:1 to 13:1.
If you don't have this information, you can set the c/s ratio to 100.

**config.yml**
Has a new attribute for sulfur deposition. The sulfur deposition is given in kg/ha annually. The deposition depends on the region and the pollution in the region. The sulfur deposition can be found in the literature or can be measured. It has been the main source of sulfur for the soil.
With new environmental regulations, the sulfur deposition decreased in the last decades, so it is important to know the sulfur deposition in the region.
It decreased from 30-40 kg/ha to 10-20 kg/ha, or less in the last decades. 
Near populated areas, the sulfur deposition can be higher, up to 50 kg/ha.

(Erfassung, Prognose und Bewertung von Stoffeinträgen und ihren Wirkungen in
Deutschland Forschungskennzahl 3707 64 200 UBA-FB 001490/ANH, 11)
Example: 
Bandenburg 2006 average SDeposition: 5.68 
Berlin 2006 average SDeposition: 7.92
Germany 2006 average SDeposition: 7.45

SDeposition: 10 // this value is properly to high, but it is just an example.

**automan.txt**
automanagemt is not implemented for sulfur yet.

**examples\parameter\FERTILIZ.TXT**

TODO: add sulfur fertilization parameters for different fertilizers.



