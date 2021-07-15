package wolfram

import (
	"encoding/xml"
	"testing"
)

const data = `
<?xml version='1.0' encoding='UTF-8'?>
<queryresult success='true'
    error='false'
    xml:space='preserve'
    numpods='4'
    datatypes=''
    timedout=''
    timedoutpods=''
    timing='2.982'
    parsetiming='0.788'
    parsetimedout='false'
    recalculate=''
    id='MSP1917208ddef655f8h83e00000h579d6aa3i0e01d'
    host='https://www4d.wolframalpha.com'
    server='22'
    related='https://www4d.wolframalpha.com/api/v1/relatedQueries.jsp?id=MSPa1918208ddef655f8h83e00003957h3d82h1gc3hi5179438809786073608'
    version='2.6'
    inputstring='integrate x^2 sin^3(x) from 0 to pi'>
  <pod title='Definite integral'
     scanner='Integral'
     id='Input'
     position='100'
     error='false'
     numsubpods='1'
     primary='true'>
    <subpod title=''>
      <img src='https://www4d.wolframalpha.com/Calculate/MSP/MSP1919208ddef655f8h83e00002003h99ai5h1i1a9?MSPStoreType=image/gif&amp;s=22'
       alt='integral_0^π x^2 sin^3(x) dx = 2/27 (9 π^2 - 40)≈3.6168'
       title='integral_0^π x^2 sin^3(x) dx = 2/27 (9 π^2 - 40)≈3.6168'
       width='284'
       height='39'
       type='Default'
       themes='1,2,3,4,5,6,7,8,9,10,11,12'
       colorinvertable='true' />
      <plaintext>integral_0^π x^2 sin^3(x) dx = 2/27 (9 π^2 - 40)≈3.6168</plaintext>
    </subpod>
    <expressiontypes count='1'>
      <expressiontype name='Default' />
    </expressiontypes>
    <states count='2'>
      <state name='More digits'
       input='Input__More digits' />
      <state name='Step-by-step solution'
       input='Input__Step-by-step solution'
       stepbystep='true' />
    </states>
  </pod>
  <pod title='Visual representation of the integral'
     scanner='Integral'
     id='VisualRepresentationOfTheIntegral'
     position='200'
     error='false'
     numsubpods='1'>
    <subpod title=''>
      <img src='https://www4d.wolframalpha.com/Calculate/MSP/MSP1920208ddef655f8h83e000039g09544f0d16g5c?MSPStoreType=image/gif&amp;s=22'
       alt='Visual representation of the integral'
       title=''
       width='220'
       height='142'
       type='2DMathPlot_1'
       themes='1,2,3,4,5,6,7,8,9,10,11,12'
       colorinvertable='true' />
      <plaintext></plaintext>
    </subpod>
    <expressiontypes count='1'>
      <expressiontype name='Default' />
    </expressiontypes>
  </pod>
  <pod title='Riemann sums'
     scanner='Integral'
     id='RiemannSums'
     position='300'
     error='false'
     numsubpods='1'>
    <subpod title=''>
      <img src='https://www4d.wolframalpha.com/Calculate/MSP/MSP1921208ddef655f8h83e00004fi7334h53heb44i?MSPStoreType=image/gif&amp;s=22'
       alt='left sum | (π (6 π^2 n^2 sin((π (n + 1))/(2 n)) csc(π/(2 n)) - 2 π^2 n^2 sin((π (n + 3))/(2 n)) csc((3 π)/(2 n)) - 3 π^2 sin((π (n - 1))/(2 n)) csc^3(π/(2 n)) + 3 π^2 sin((3 π (n - 1))/(2 n)) csc^3(π/(2 n)) - π^2 sin((3 π (n - 3))/(2 n)) csc^3((3 π)/(2 n)) - π^2 sin((3 π (n + 1))/(2 n)) csc^3((3 π)/(2 n)) - 9 π^2 sin(π/n) csc^2(π/(2 n)) + 3 π^2 sin((3 π)/n) csc^2((3 π)/(2 n)) + 6 π^2 sin((π (n - 1))/(2 n)) csc(π/(2 n)) - 2 π^2 sin((π (n - 3))/(2 n)) csc((3 π)/(2 n))))/(16 n^3) = 2/27 (9 π^2 - 40) + O((1/n)^4)
(assuming subintervals of equal length)'
       title='left sum | (π (6 π^2 n^2 sin((π (n + 1))/(2 n)) csc(π/(2 n)) - 2 π^2 n^2 sin((π (n + 3))/(2 n)) csc((3 π)/(2 n)) - 3 π^2 sin((π (n - 1))/(2 n)) csc^3(π/(2 n)) + 3 π^2 sin((3 π (n - 1))/(2 n)) csc^3(π/(2 n)) - π^2 sin((3 π (n - 3))/(2 n)) csc^3((3 π)/(2 n)) - π^2 sin((3 π (n + 1))/(2 n)) csc^3((3 π)/(2 n)) - 9 π^2 sin(π/n) csc^2(π/(2 n)) + 3 π^2 sin((3 π)/n) csc^2((3 π)/(2 n)) + 6 π^2 sin((π (n - 1))/(2 n)) csc(π/(2 n)) - 2 π^2 sin((π (n - 3))/(2 n)) csc((3 π)/(2 n))))/(16 n^3) = 2/27 (9 π^2 - 40) + O((1/n)^4)
(assuming subintervals of equal length)'
       width='546'
       height='179'
       type='Grid'
       themes='1,2,3,4,5,6,7,8,9,10,11,12'
       colorinvertable='true' />
      <plaintext>left sum | (π (6 π^2 n^2 sin((π (n + 1))/(2 n)) csc(π/(2 n)) - 2 π^2 n^2 sin((π (n + 3))/(2 n)) csc((3 π)/(2 n)) - 3 π^2 sin((π (n - 1))/(2 n)) csc^3(π/(2 n)) + 3 π^2 sin((3 π (n - 1))/(2 n)) csc^3(π/(2 n)) - π^2 sin((3 π (n - 3))/(2 n)) csc^3((3 π)/(2 n)) - π^2 sin((3 π (n + 1))/(2 n)) csc^3((3 π)/(2 n)) - 9 π^2 sin(π/n) csc^2(π/(2 n)) + 3 π^2 sin((3 π)/n) csc^2((3 π)/(2 n)) + 6 π^2 sin((π (n - 1))/(2 n)) csc(π/(2 n)) - 2 π^2 sin((π (n - 3))/(2 n)) csc((3 π)/(2 n))))/(16 n^3) = 2/27 (9 π^2 - 40) + O((1/n)^4)
(assuming subintervals of equal length)</plaintext>
    </subpod>
    <expressiontypes count='1'>
      <expressiontype name='Grid' />
    </expressiontypes>
    <states count='1'>
      <state name='More cases'
       input='RiemannSums__More cases' />
    </states>
    <infos count='1'>
      <info text='csc(x) is the cosecant function'>
        <img src='https://www4d.wolframalpha.com/Calculate/MSP/MSP1922208ddef655f8h83e00001d8d17d2c0gi92ca?MSPStoreType=image/gif&amp;s=22'
        alt='csc(x) is the cosecant function'
        title='csc(x) is the cosecant function'
        width='198'
        height='19' />
        <link url='http://reference.wolfram.com/language/ref/Csc.html'
        text='Documentation'
        title='Mathematica' />
        <link url='http://functions.wolfram.com/ElementaryFunctions/Csc'
        text='Properties'
        title='Wolfram Functions Site' />
        <link url='http://mathworld.wolfram.com/Cosecant.html'
        text='Definition'
        title='MathWorld' />
      </info>
    </infos>
  </pod>
  <pod title='Indefinite integral'
     scanner='Integral'
     id='IndefiniteIntegral'
     position='400'
     error='false'
     numsubpods='1'>
    <subpod title=''>
      <img src='https://www4d.wolframalpha.com/Calculate/MSP/MSP1923208ddef655f8h83e00002gif413iab4i2956?MSPStoreType=image/gif&amp;s=22'
       alt='integral x^2 sin^3(x) dx = 1/108 (-81 (x^2 - 2) cos(x) + (9 x^2 - 2) cos(3 x) - 6 x (sin(3 x) - 27 sin(x))) + constant'
       title='integral x^2 sin^3(x) dx = 1/108 (-81 (x^2 - 2) cos(x) + (9 x^2 - 2) cos(3 x) - 6 x (sin(3 x) - 27 sin(x))) + constant'
       width='493'
       height='92'
       type='Default'
       themes='1,2,3,4,5,6,7,8,9,10,11,12'
       colorinvertable='true' />
      <plaintext>integral x^2 sin^3(x) dx = 1/108 (-81 (x^2 - 2) cos(x) + (9 x^2 - 2) cos(3 x) - 6 x (sin(3 x) - 27 sin(x))) + constant</plaintext>
    </subpod>
    <expressiontypes count='1'>
      <expressiontype name='Default' />
    </expressiontypes>
    <states count='1'>
      <state name='Step-by-step solution'
       input='IndefiniteIntegral__Step-by-step solution'
       stepbystep='true' />
    </states>
  </pod>
</queryresult>
`

const data2 = `
<?xml version='1.0' encoding='UTF-8'?>
<queryresult success='true'
    error='false'
    xml:space='preserve'
    numpods='2'
    datatypes='MusicWork'
    timedout=''
    timedoutpods=''
    timing='1.674'
    parsetiming='0.74'
    parsetimedout='false'
    recalculate=''
    id='MSP480917222eiab3bcd712000049a9g06e0f48c66i'
    host='https://www4f.wolframalpha.com'
    server='10'
    related='https://www4f.wolframalpha.com/api/v1/relatedQueries.jsp?id=MSPa481017222eiab3bcd71200005i40cif6ia2040408426677913092218004'
    version='2.6'
    inputstring='i don&apos;t want to live'>
  <pod title='Input interpretation'
     scanner='Identity'
     id='Input'
     position='100'
     error='false'
     numsubpods='1'>
    <subpod title=''>
      <img src='https://www4f.wolframalpha.com/Calculate/MSP/MSP481117222eiab3bcd71200003470beid3hbc79i8?MSPStoreType=image/gif&amp;s=10'
       alt='I Don&apos;t Wanna Live (music work)'
       title='I Don&apos;t Wanna Live (music work)'
       width='222'
       height='19'
       type='Default'
       themes='1,2,3,4,5,6,7,8,9,10,11,12'
       colorinvertable='true' />
      <plaintext>I Don&apos;t Wanna Live (music work)</plaintext>
    </subpod>
    <expressiontypes count='1'>
      <expressiontype name='Default' />
    </expressiontypes>
  </pod>
  <pod title='Recordings'
     scanner='Data'
     id='Recordings:MusicWorkData'
     position='200'
     error='false'
     numsubpods='1'>
    <subpod title=''>
      <microsources>
        <microsource>MusicWorkData</microsource>
      </microsources>
      <datasources>
        <datasource>MusicBrainz</datasource>
      </datasources>
      <img src='https://www4f.wolframalpha.com/Calculate/MSP/MSP481217222eiab3bcd71200000g1gdi2h67e8aa91?MSPStoreType=image/gif&amp;s=10'
       alt='album | music act | release date
My Life&apos;s Been a Country Song | Chris Cagle | Tuesday, February 19, 2008'
       title='album | music act | release date
My Life&apos;s Been a Country Song | Chris Cagle | Tuesday, February 19, 2008'
       width='526'
       height='66'
       type='Grid'
       themes='1,2,3,4,5,6,7,8,9,10,11,12'
       colorinvertable='true' />
      <plaintext>album | music act | release date
My Life&apos;s Been a Country Song | Chris Cagle | Tuesday, February 19, 2008</plaintext>
    </subpod>
    <expressiontypes count='1'>
      <expressiontype name='Grid' />
    </expressiontypes>
  </pod>
  <assumptions count='1'>
    <assumption type='SubCategory'
      word='i don&apos;t want to live'
      template='Assuming ${desc1}. Use ${desc2} instead'
      count='2'>
      <value name='IDontWannaLive::x96n6'
       desc='I Don&apos;t Wanna Live (Chris Cagle)'
       input='*DPClash.MusicWorkE.i+don%27t+want+to+live-_*IDontWannaLive%3A%3Ax96n6-' />
      <value name='IDontWanttoLive::y87c6'
       desc='I Don&apos;t Want to Live (Josh Gracin)'
       input='*DPClash.MusicWorkE.i+don%27t+want+to+live-_*IDontWanttoLive%3A%3Ay87c6-' />
    </assumption>
  </assumptions>
  <sources count='1'>
    <source url='https://www4f.wolframalpha.com/sources/MusicWorkDataSourceInformationNotes.html'
      text='Music work data' />
  </sources>
</queryresult>`

func TestXMLParsing(t *testing.T) {
	var result QueryResult
	if err := xml.Unmarshal([]byte(data2), &result); err != nil {
		t.Error(err)
	}

	t.Log(result)
}
