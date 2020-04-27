package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const SAMPLE = `[{"locale": "en", "description": "PMC\nPMC\n52\n1364\n586\nTRIAD\nPANTHER\n2003\n2182\n2686\nWalter Scott THE TALISMAN\nISBN 0 14\n00.5826 5\n88\nInge Scholl Die Weiße Rose\nwide sarg asso sea jean rhyS\nCOLETTE ROSSANT\nAPRICOTS ON THE NILE\nJohn Cowper Powys Wolf Solent\nGeorge Orwell Down end Out in Paris and London\n5ON O14\n00.0838 1\nGeorge Orwell Animal Farm\nFlann O'Brrien\nAt Swim-Two-Birds\nWOODE\nS\nNORWYGIAN WOUD\nENORWEGIAN WOOD\nrakami\nlarukt\nIris Murdoch A Severed Head\niris Murdoch An Unofficial Rose\nISBN 0 14\n00.2154 X\nIris Murdoch The Time of the Angels\n14 002848 X\nIRIS MURDOCH\nTHE SEA,THE SEA\nTropic of Capricorn Henry Miller\nThe Screwtape Letters C.S. Lewis\nAlan Jay Lerner\nMy Fair Lady\nD-H.LAWRENCE Lady Chatterley's Lover\nISRN 0 14\nGARRISON KEILLOR Radio Romance\nFranz Kafka America\nISBN 0 14\n00.2639 8\nFranz Kafka The Trial\nISBN 0 14\n00.0907 8\n", "boundingPoly": {"vertices": [{"x": 97, "y": 660}, {"x": 3865, "y": 660}, {"x": 3865, "y": 2667}, {"x": 97, "y": 2667}]}}, {"description": "PMC", "boundingPoly": {"vertices": [{"x": 333, "y": 688}, {"x": 423, "y": 690}, {"x": 422, "y": 726}, {"x": 332, "y": 724}]}}, {"description": "PMC", "boundingPoly": {"vertices": [{"x": 146, "y": 741}, {"x": 241, "y": 745}, {"x": 239, "y": 780}, {"x": 145, "y": 776}]}}, {"description": "52", "boundingPoly": {"vertices": [{"x": 2530, "y": 1136}, {"x": 2566, "y": 1136}, {"x": 2566, "y": 1166}, {"x": 2530, "y": 1166}]}}, {"description": "1364", "boundingPoly": {"vertices": [{"x": 990, "y": 2418}, {"x": 1053, "y": 2417}, {"x": 1053, "y": 2442}, {"x": 990, "y": 2443}]}}, {"description": "586", "boundingPoly": {"vertices": [{"x": 3679, "y": 2475}, {"x": 3727, "y": 2474}, {"x": 3728, "y": 2504}, {"x": 3680, "y": 2505}]}}, {"description": "TRIAD", "boundingPoly": {"vertices": [{"x": 1605, "y": 2532}, {"x": 1701, "y": 2532}, {"x": 1701, "y": 2561}, {"x": 1605, "y": 2561}]}}, {"description": "PANTHER", "boundingPoly": {"vertices": [{"x": 1573, "y": 2570}, {"x": 1727, "y": 2569}, {"x": 1727, "y": 2600}, {"x": 1573, "y": 2601}]}}, {"description": "2003", "boundingPoly": {"vertices": [{"x": 2254, "y": 2594}, {"x": 2312, "y": 2591}, {"x": 2313, "y": 2620}, {"x": 2256, "y": 2623}]}}, {"description": "2182", "boundingPoly": {"vertices": [{"x": 3169, "y": 2601}, {"x": 3225, "y": 2598}, {"x": 3226, "y": 2627}, {"x": 3170, "y": 2630}]}}, {"description": "2686", "boundingPoly": {"vertices": [{"x": 2689, "y": 2623}, {"x": 2748, "y": 2621}, {"x": 2749, "y": 2649}, {"x": 2690, "y": 2651}]}}, {"description": "Walter", "boundingPoly": {"vertices": [{"x": 3848, "y": 692}, {"x": 3849, "y": 899}, {"x": 3798, "y": 899}, {"x": 3797, "y": 692}]}}, {"description": "Scott", "boundingPoly": {"vertices": [{"x": 3849, "y": 914}, {"x": 3850, "y": 1075}, {"x": 3797, "y": 1075}, {"x": 3796, "y": 914}]}}, {"description": "THE", "boundingPoly": {"vertices": [{"x": 3849, "y": 1114}, {"x": 3850, "y": 1255}, {"x": 3798, "y": 1255}, {"x": 3797, "y": 1114}]}}, {"description": "TALISMAN", "boundingPoly": {"vertices": [{"x": 3850, "y": 1274}, {"x": 3852, "y": 1631}, {"x": 3802, "y": 1631}, {"x": 3800, "y": 1274}]}}, {"description": "ISBN", "boundingPoly": {"vertices": [{"x": 3864, "y": 2366}, {"x": 3865, "y": 2423}, {"x": 3844, "y": 2423}, {"x": 3843, "y": 2366}]}}, {"description": "0", "boundingPoly": {"vertices": [{"x": 3865, "y": 2437}, {"x": 3865, "y": 2450}, {"x": 3844, "y": 2450}, {"x": 3844, "y": 2437}]}}, {"description": "14", "boundingPoly": {"vertices": [{"x": 3865, "y": 2465}, {"x": 3865, "y": 2490}, {"x": 3845, "y": 2490}, {"x": 3845, "y": 2465}]}}, {"description": "00.5826", "boundingPoly": {"vertices": [{"x": 3835, "y": 2364}, {"x": 3837, "y": 2465}, {"x": 3815, "y": 2465}, {"x": 3813, "y": 2364}]}}, {"description": "5", "boundingPoly": {"vertices": [{"x": 3836, "y": 2478}, {"x": 3836, "y": 2490}, {"x": 3815, "y": 2490}, {"x": 3815, "y": 2478}]}}, {"description": "88", "boundingPoly": {"vertices": [{"x": 3657, "y": 756}, {"x": 3656, "y": 700}, {"x": 3698, "y": 699}, {"x": 3699, "y": 755}]}}, {"description": "Inge", "boundingPoly": {"vertices": [{"x": 3665, "y": 1597}, {"x": 3663, "y": 1487}, {"x": 3711, "y": 1486}, {"x": 3713, "y": 1596}]}}, {"description": "Scholl", "boundingPoly": {"vertices": [{"x": 3664, "y": 1466}, {"x": 3662, "y": 1307}, {"x": 3704, "y": 1306}, {"x": 3706, "y": 1465}]}}, {"description": "Die", "boundingPoly": {"vertices": [{"x": 3660, "y": 1248}, {"x": 3659, "y": 1164}, {"x": 3703, "y": 1163}, {"x": 3704, "y": 1247}]}}, {"description": "Weiße", "boundingPoly": {"vertices": [{"x": 3658, "y": 1150}, {"x": 3656, "y": 991}, {"x": 3702, "y": 990}, {"x": 3704, "y": 1149}]}}, {"description": "Rose", "boundingPoly": {"vertices": [{"x": 3656, "y": 979}, {"x": 3654, "y": 846}, {"x": 3702, "y": 845}, {"x": 3704, "y": 978}]}}, {"description": "wide", "boundingPoly": {"vertices": [{"x": 3552, "y": 1144}, {"x": 3557, "y": 1287}, {"x": 3520, "y": 1288}, {"x": 3515, "y": 1145}]}}, {"description": "sarg", "boundingPoly": {"vertices": [{"x": 3559, "y": 1327}, {"x": 3563, "y": 1449}, {"x": 3514, "y": 1451}, {"x": 3510, "y": 1329}]}}, {"description": "asso", "boundingPoly": {"vertices": [{"x": 3564, "y": 1467}, {"x": 3569, "y": 1599}, {"x": 3519, "y": 1601}, {"x": 3514, "y": 1469}]}}, {"description": "sea", "boundingPoly": {"vertices": [{"x": 3557, "y": 1628}, {"x": 3560, "y": 1723}, {"x": 3533, "y": 1724}, {"x": 3530, "y": 1629}]}}, {"description": "jean", "boundingPoly": {"vertices": [{"x": 3575, "y": 1784}, {"x": 3580, "y": 1927}, {"x": 3533, "y": 1929}, {"x": 3528, "y": 1786}]}}, {"description": "rhyS", "boundingPoly": {"vertices": [{"x": 3581, "y": 1957}, {"x": 3586, "y": 2092}, {"x": 3537, "y": 2094}, {"x": 3532, "y": 1959}]}}, {"description": "COLETTE", "boundingPoly": {"vertices": [{"x": 3426, "y": 1893}, {"x": 3429, "y": 2111}, {"x": 3401, "y": 2111}, {"x": 3398, "y": 1893}]}}, {"description": "ROSSANT", "boundingPoly": {"vertices": [{"x": 3431, "y": 2136}, {"x": 3434, "y": 2345}, {"x": 3405, "y": 2345}, {"x": 3402, "y": 2136}]}}, {"description": "APRICOTS", "boundingPoly": {"vertices": [{"x": 3412, "y": 789}, {"x": 3418, "y": 1055}, {"x": 3380, "y": 1056}, {"x": 3374, "y": 790}]}}, {"description": "ON", "boundingPoly": {"vertices": [{"x": 3419, "y": 1123}, {"x": 3420, "y": 1173}, {"x": 3375, "y": 1174}, {"x": 3374, "y": 1124}]}}, {"description": "THE", "boundingPoly": {"vertices": [{"x": 3412, "y": 1194}, {"x": 3414, "y": 1295}, {"x": 3385, "y": 1296}, {"x": 3383, "y": 1195}]}}, {"description": "NILE", "boundingPoly": {"vertices": [{"x": 3424, "y": 1330}, {"x": 3427, "y": 1459}, {"x": 3389, "y": 1460}, {"x": 3386, "y": 1331}]}}, {"description": "John", "boundingPoly": {"vertices": [{"x": 3191, "y": 1120}, {"x": 3194, "y": 1215}, {"x": 3153, "y": 1216}, {"x": 3150, "y": 1121}]}}, {"description": "Cowper", "boundingPoly": {"vertices": [{"x": 3191, "y": 1240}, {"x": 3195, "y": 1399}, {"x": 3150, "y": 1400}, {"x": 3146, "y": 1241}]}}, {"description": "Powys", "boundingPoly": {"vertices": [{"x": 3194, "y": 1420}, {"x": 3197, "y": 1548}, {"x": 3154, "y": 1549}, {"x": 3151, "y": 1421}]}}, {"description": "Wolf", "boundingPoly": {"vertices": [{"x": 3200, "y": 1606}, {"x": 3203, "y": 1714}, {"x": 3169, "y": 1715}, {"x": 3166, "y": 1607}]}}, {"description": "Solent", "boundingPoly": {"vertices": [{"x": 3202, "y": 1731}, {"x": 3205, "y": 1853}, {"x": 3170, "y": 1854}, {"x": 3167, "y": 1732}]}}, {"description": "George", "boundingPoly": {"vertices": [{"x": 2977, "y": 1233}, {"x": 2979, "y": 1377}, {"x": 2938, "y": 1378}, {"x": 2936, "y": 1234}]}}, {"description": "Orwell", "boundingPoly": {"vertices": [{"x": 2979, "y": 1401}, {"x": 2981, "y": 1538}, {"x": 2940, "y": 1539}, {"x": 2938, "y": 1402}]}}, {"description": "Down", "boundingPoly": {"vertices": [{"x": 2981, "y": 1569}, {"x": 2983, "y": 1685}, {"x": 2943, "y": 1686}, {"x": 2941, "y": 1570}]}}, {"description": "end", "boundingPoly": {"vertices": [{"x": 2980, "y": 1701}, {"x": 2981, "y": 1777}, {"x": 2950, "y": 1777}, {"x": 2949, "y": 1701}]}}, {"description": "Out", "boundingPoly": {"vertices": [{"x": 2983, "y": 1797}, {"x": 2984, "y": 1873}, {"x": 2951, "y": 1873}, {"x": 2950, "y": 1797}]}}, {"description": "in", "boundingPoly": {"vertices": [{"x": 2986, "y": 1884}, {"x": 2987, "y": 1930}, {"x": 2946, "y": 1931}, {"x": 2945, "y": 1885}]}}, {"description": "Paris", "boundingPoly": {"vertices": [{"x": 2982, "y": 1945}, {"x": 2984, "y": 2059}, {"x": 2956, "y": 2059}, {"x": 2954, "y": 1945}]}}, {"description": "and", "boundingPoly": {"vertices": [{"x": 2989, "y": 2073}, {"x": 2990, "y": 2147}, {"x": 2949, "y": 2148}, {"x": 2948, "y": 2074}]}}, {"description": "London", "boundingPoly": {"vertices": [{"x": 2990, "y": 2164}, {"x": 2992, "y": 2322}, {"x": 2951, "y": 2323}, {"x": 2949, "y": 2165}]}}, {"description": "5ON", "boundingPoly": {"vertices": [{"x": 2876, "y": 2312}, {"x": 2878, "y": 2412}, {"x": 2849, "y": 2413}, {"x": 2847, "y": 2313}]}}, {"description": "O14", "boundingPoly": {"vertices": [{"x": 2879, "y": 2424}, {"x": 2880, "y": 2474}, {"x": 2857, "y": 2474}, {"x": 2856, "y": 2425}]}}, {"description": "00.0838", "boundingPoly": {"vertices": [{"x": 2848, "y": 2345}, {"x": 2850, "y": 2453}, {"x": 2823, "y": 2453}, {"x": 2821, "y": 2345}]}}, {"description": "1", "boundingPoly": {"vertices": [{"x": 2846, "y": 2465}, {"x": 2846, "y": 2473}, {"x": 2825, "y": 2473}, {"x": 2825, "y": 2465}]}}, {"description": "George", "boundingPoly": {"vertices": [{"x": 2857, "y": 712}, {"x": 2862, "y": 935}, {"x": 2790, "y": 937}, {"x": 2785, "y": 714}]}}, {"description": "Orwell", "boundingPoly": {"vertices": [{"x": 2862, "y": 948}, {"x": 2867, "y": 1151}, {"x": 2811, "y": 1152}, {"x": 2806, "y": 949}]}}, {"description": "Animal", "boundingPoly": {"vertices": [{"x": 2864, "y": 1204}, {"x": 2869, "y": 1427}, {"x": 2814, "y": 1428}, {"x": 2809, "y": 1205}]}}, {"description": "Farm", "boundingPoly": {"vertices": [{"x": 2866, "y": 1448}, {"x": 2870, "y": 1603}, {"x": 2815, "y": 1604}, {"x": 2811, "y": 1449}]}}, {"description": "Flann", "boundingPoly": {"vertices": [{"x": 2741, "y": 688}, {"x": 2741, "y": 829}, {"x": 2705, "y": 829}, {"x": 2705, "y": 688}]}}, {"description": "O'Brrien", "boundingPoly": {"vertices": [{"x": 2737, "y": 837}, {"x": 2737, "y": 991}, {"x": 2707, "y": 991}, {"x": 2707, "y": 837}]}}, {"description": "At", "boundingPoly": {"vertices": [{"x": 2700, "y": 720}, {"x": 2700, "y": 757}, {"x": 2673, "y": 757}, {"x": 2673, "y": 720}]}}, {"description": "Swim-Two-Birds", "boundingPoly": {"vertices": [{"x": 2703, "y": 772}, {"x": 2703, "y": 1076}, {"x": 2668, "y": 1076}, {"x": 2668, "y": 772}]}}, {"description": "WOODE", "boundingPoly": {"vertices": [{"x": 2588, "y": 1646}, {"x": 2583, "y": 1915}, {"x": 2552, "y": 1914}, {"x": 2557, "y": 1645}]}}, {"description": "S", "boundingPoly": {"vertices": [{"x": 2579, "y": 1147}, {"x": 2580, "y": 1223}, {"x": 2515, "y": 1224}, {"x": 2514, "y": 1148}]}}, {"description": "NORWYGIAN", "boundingPoly": {"vertices": [{"x": 2583, "y": 1248}, {"x": 2588, "y": 1630}, {"x": 2515, "y": 1631}, {"x": 2510, "y": 1249}]}}, {"description": "WOUD", "boundingPoly": {"vertices": [{"x": 2588, "y": 1635}, {"x": 2591, "y": 1861}, {"x": 2518, "y": 1862}, {"x": 2515, "y": 1636}]}}, {"description": "ENORWEGIAN", "boundingPoly": {"vertices": [{"x": 2447, "y": 1142}, {"x": 2451, "y": 1629}, {"x": 2383, "y": 1630}, {"x": 2379, "y": 1143}]}}, {"description": "WOOD", "boundingPoly": {"vertices": [{"x": 2451, "y": 1656}, {"x": 2453, "y": 1867}, {"x": 2385, "y": 1868}, {"x": 2383, "y": 1657}]}}, {"description": "rakami", "boundingPoly": {"vertices": [{"x": 2439, "y": 2144}, {"x": 2441, "y": 2307}, {"x": 2408, "y": 2307}, {"x": 2406, "y": 2144}]}}, {"description": "larukt", "boundingPoly": {"vertices": [{"x": 2433, "y": 1974}, {"x": 2433, "y": 2090}, {"x": 2432, "y": 2090}, {"x": 2432, "y": 1974}]}}, {"description": "Iris", "boundingPoly": {"vertices": [{"x": 2308, "y": 841}, {"x": 2307, "y": 923}, {"x": 2261, "y": 923}, {"x": 2262, "y": 841}]}}, {"description": "Murdoch", "boundingPoly": {"vertices": [{"x": 2308, "y": 943}, {"x": 2306, "y": 1147}, {"x": 2260, "y": 1147}, {"x": 2262, "y": 943}]}}, {"description": "A", "boundingPoly": {"vertices": [{"x": 2299, "y": 1180}, {"x": 2299, "y": 1209}, {"x": 2266, "y": 1209}, {"x": 2266, "y": 1180}]}}, {"description": "Severed", "boundingPoly": {"vertices": [{"x": 2299, "y": 1232}, {"x": 2297, "y": 1425}, {"x": 2262, "y": 1425}, {"x": 2264, "y": 1232}]}}, {"description": "Head", "boundingPoly": {"vertices": [{"x": 2299, "y": 1452}, {"x": 2298, "y": 1573}, {"x": 2265, "y": 1573}, {"x": 2266, "y": 1452}]}}, {"description": "iris", "boundingPoly": {"vertices": [{"x": 2164, "y": 714}, {"x": 2163, "y": 795}, {"x": 2103, "y": 794}, {"x": 2104, "y": 713}]}}, {"description": "Murdoch", "boundingPoly": {"vertices": [{"x": 2163, "y": 818}, {"x": 2160, "y": 1045}, {"x": 2100, "y": 1044}, {"x": 2103, "y": 817}]}}, {"description": "An", "boundingPoly": {"vertices": [{"x": 2154, "y": 1070}, {"x": 2153, "y": 1137}, {"x": 2111, "y": 1137}, {"x": 2112, "y": 1070}]}}, {"description": "Unofficial", "boundingPoly": {"vertices": [{"x": 2153, "y": 1158}, {"x": 2150, "y": 1393}, {"x": 2107, "y": 1392}, {"x": 2110, "y": 1157}]}}, {"description": "Rose", "boundingPoly": {"vertices": [{"x": 2150, "y": 1414}, {"x": 2148, "y": 1545}, {"x": 2106, "y": 1544}, {"x": 2108, "y": 1414}]}}, {"description": "ISBN", "boundingPoly": {"vertices": [{"x": 2150, "y": 2418}, {"x": 2149, "y": 2481}, {"x": 2126, "y": 2480}, {"x": 2127, "y": 2417}]}}, {"description": "0", "boundingPoly": {"vertices": [{"x": 2154, "y": 2498}, {"x": 2154, "y": 2518}, {"x": 2120, "y": 2517}, {"x": 2120, "y": 2497}]}}, {"description": "14", "boundingPoly": {"vertices": [{"x": 2148, "y": 2516}, {"x": 2147, "y": 2546}, {"x": 2124, "y": 2545}, {"x": 2125, "y": 2515}]}}, {"description": "00.2154", "boundingPoly": {"vertices": [{"x": 2118, "y": 2417}, {"x": 2118, "y": 2521}, {"x": 2094, "y": 2521}, {"x": 2094, "y": 2417}]}}, {"description": "X", "boundingPoly": {"vertices": [{"x": 2117, "y": 2529}, {"x": 2117, "y": 2545}, {"x": 2094, "y": 2545}, {"x": 2094, "y": 2529}]}}, {"description": "Iris", "boundingPoly": {"vertices": [{"x": 1962, "y": 660}, {"x": 1962, "y": 755}, {"x": 1901, "y": 755}, {"x": 1901, "y": 660}]}}, {"description": "Murdoch", "boundingPoly": {"vertices": [{"x": 1962, "y": 768}, {"x": 1962, "y": 1003}, {"x": 1901, "y": 1003}, {"x": 1901, "y": 768}]}}, {"description": "The", "boundingPoly": {"vertices": [{"x": 1954, "y": 1028}, {"x": 1954, "y": 1129}, {"x": 1911, "y": 1129}, {"x": 1911, "y": 1028}]}}, {"description": "Time", "boundingPoly": {"vertices": [{"x": 1956, "y": 1144}, {"x": 1956, "y": 1273}, {"x": 1911, "y": 1273}, {"x": 1911, "y": 1144}]}}, {"description": "of", "boundingPoly": {"vertices": [{"x": 1956, "y": 1290}, {"x": 1956, "y": 1339}, {"x": 1911, "y": 1339}, {"x": 1911, "y": 1290}]}}, {"description": "the", "boundingPoly": {"vertices": [{"x": 1956, "y": 1354}, {"x": 1956, "y": 1435}, {"x": 1913, "y": 1435}, {"x": 1913, "y": 1354}]}}, {"description": "Angels", "boundingPoly": {"vertices": [{"x": 1956, "y": 1452}, {"x": 1956, "y": 1631}, {"x": 1903, "y": 1631}, {"x": 1903, "y": 1452}]}}, {"description": "14", "boundingPoly": {"vertices": [{"x": 1942, "y": 2358}, {"x": 1942, "y": 2386}, {"x": 1920, "y": 2386}, {"x": 1920, "y": 2358}]}}, {"description": "002848", "boundingPoly": {"vertices": [{"x": 1943, "y": 2391}, {"x": 1943, "y": 2478}, {"x": 1919, "y": 2478}, {"x": 1919, "y": 2391}]}}, {"description": "X", "boundingPoly": {"vertices": [{"x": 1944, "y": 2484}, {"x": 1944, "y": 2500}, {"x": 1921, "y": 2500}, {"x": 1921, "y": 2484}]}}, {"description": "IRIS", "boundingPoly": {"vertices": [{"x": 1800, "y": 682}, {"x": 1800, "y": 931}, {"x": 1673, "y": 931}, {"x": 1673, "y": 682}]}}, {"description": "MURDOCH", "boundingPoly": {"vertices": [{"x": 1800, "y": 958}, {"x": 1800, "y": 1443}, {"x": 1673, "y": 1443}, {"x": 1673, "y": 958}]}}, {"description": "THE", "boundingPoly": {"vertices": [{"x": 1672, "y": 688}, {"x": 1671, "y": 907}, {"x": 1552, "y": 907}, {"x": 1553, "y": 688}]}}, {"description": "SEA,THE", "boundingPoly": {"vertices": [{"x": 1664, "y": 924}, {"x": 1663, "y": 1343}, {"x": 1539, "y": 1343}, {"x": 1540, "y": 924}]}}, {"description": "SEA", "boundingPoly": {"vertices": [{"x": 1662, "y": 1364}, {"x": 1661, "y": 1555}, {"x": 1549, "y": 1555}, {"x": 1550, "y": 1364}]}}, {"description": "Tropic", "boundingPoly": {"vertices": [{"x": 1394, "y": 1022}, {"x": 1393, "y": 1229}, {"x": 1319, "y": 1229}, {"x": 1320, "y": 1022}]}}, {"description": "of", "boundingPoly": {"vertices": [{"x": 1401, "y": 1260}, {"x": 1401, "y": 1319}, {"x": 1334, "y": 1319}, {"x": 1334, "y": 1260}]}}, {"description": "Capricorn", "boundingPoly": {"vertices": [{"x": 1393, "y": 1356}, {"x": 1392, "y": 1691}, {"x": 1316, "y": 1691}, {"x": 1317, "y": 1356}]}}, {"description": "Henry", "boundingPoly": {"vertices": [{"x": 1389, "y": 1770}, {"x": 1388, "y": 1969}, {"x": 1313, "y": 1969}, {"x": 1314, "y": 1770}]}}, {"description": "Miller", "boundingPoly": {"vertices": [{"x": 1391, "y": 1998}, {"x": 1390, "y": 2191}, {"x": 1310, "y": 2191}, {"x": 1311, "y": 1998}]}}, {"description": "The", "boundingPoly": {"vertices": [{"x": 1220, "y": 858}, {"x": 1220, "y": 959}, {"x": 1175, "y": 959}, {"x": 1175, "y": 858}]}}, {"description": "Screwtape", "boundingPoly": {"vertices": [{"x": 1218, "y": 978}, {"x": 1218, "y": 1237}, {"x": 1164, "y": 1237}, {"x": 1164, "y": 978}]}}, {"description": "Letters", "boundingPoly": {"vertices": [{"x": 1215, "y": 1256}, {"x": 1215, "y": 1431}, {"x": 1174, "y": 1431}, {"x": 1174, "y": 1256}]}}, {"description": "C.S.", "boundingPoly": {"vertices": [{"x": 1217, "y": 1490}, {"x": 1217, "y": 1583}, {"x": 1174, "y": 1583}, {"x": 1174, "y": 1490}]}}, {"description": "Lewis", "boundingPoly": {"vertices": [{"x": 1215, "y": 1606}, {"x": 1215, "y": 1733}, {"x": 1172, "y": 1733}, {"x": 1172, "y": 1606}]}}, {"description": "Alan", "boundingPoly": {"vertices": [{"x": 1064, "y": 1196}, {"x": 1063, "y": 1281}, {"x": 1027, "y": 1281}, {"x": 1028, "y": 1196}]}}, {"description": "Jay", "boundingPoly": {"vertices": [{"x": 1063, "y": 1302}, {"x": 1062, "y": 1361}, {"x": 1016, "y": 1360}, {"x": 1017, "y": 1301}]}}, {"description": "Lerner", "boundingPoly": {"vertices": [{"x": 1062, "y": 1382}, {"x": 1060, "y": 1511}, {"x": 1024, "y": 1511}, {"x": 1026, "y": 1382}]}}, {"description": "My", "boundingPoly": {"vertices": [{"x": 1054, "y": 1774}, {"x": 1053, "y": 1831}, {"x": 1005, "y": 1830}, {"x": 1006, "y": 1773}]}}, {"description": "Fair", "boundingPoly": {"vertices": [{"x": 1055, "y": 1852}, {"x": 1054, "y": 1923}, {"x": 1017, "y": 1923}, {"x": 1018, "y": 1852}]}}, {"description": "Lady", "boundingPoly": {"vertices": [{"x": 1052, "y": 1944}, {"x": 1051, "y": 2033}, {"x": 1002, "y": 2032}, {"x": 1003, "y": 1943}]}}, {"description": "D-H.LAWRENCE", "boundingPoly": {"vertices": [{"x": 924, "y": 760}, {"x": 912, "y": 1511}, {"x": 848, "y": 1510}, {"x": 860, "y": 759}]}}, {"description": "Lady", "boundingPoly": {"vertices": [{"x": 923, "y": 1570}, {"x": 920, "y": 1731}, {"x": 844, "y": 1730}, {"x": 847, "y": 1569}]}}, {"description": "Chatterley's", "boundingPoly": {"vertices": [{"x": 920, "y": 1762}, {"x": 914, "y": 2123}, {"x": 838, "y": 2122}, {"x": 844, "y": 1761}]}}, {"description": "Lover", "boundingPoly": {"vertices": [{"x": 898, "y": 2136}, {"x": 895, "y": 2323}, {"x": 840, "y": 2322}, {"x": 843, "y": 2135}]}}, {"description": "ISRN", "boundingPoly": {"vertices": [{"x": 894, "y": 2381}, {"x": 895, "y": 2439}, {"x": 874, "y": 2439}, {"x": 873, "y": 2381}]}}, {"description": "0", "boundingPoly": {"vertices": [{"x": 895, "y": 2452}, {"x": 895, "y": 2465}, {"x": 875, "y": 2465}, {"x": 875, "y": 2452}]}}, {"description": "14", "boundingPoly": {"vertices": [{"x": 895, "y": 2481}, {"x": 896, "y": 2507}, {"x": 876, "y": 2507}, {"x": 875, "y": 2481}]}}, {"description": "GARRISON", "boundingPoly": {"vertices": [{"x": 654, "y": 788}, {"x": 647, "y": 1163}, {"x": 592, "y": 1162}, {"x": 599, "y": 787}]}}, {"description": "KEILLOR", "boundingPoly": {"vertices": [{"x": 642, "y": 1192}, {"x": 636, "y": 1499}, {"x": 585, "y": 1498}, {"x": 591, "y": 1191}]}}, {"description": "Radio", "boundingPoly": {"vertices": [{"x": 637, "y": 1574}, {"x": 634, "y": 1750}, {"x": 579, "y": 1749}, {"x": 582, "y": 1573}]}}, {"description": "Romance", "boundingPoly": {"vertices": [{"x": 639, "y": 1780}, {"x": 633, "y": 2068}, {"x": 573, "y": 2067}, {"x": 579, "y": 1779}]}}, {"description": "Franz", "boundingPoly": {"vertices": [{"x": 396, "y": 1022}, {"x": 394, "y": 1191}, {"x": 337, "y": 1190}, {"x": 339, "y": 1021}]}}, {"description": "Kafka", "boundingPoly": {"vertices": [{"x": 389, "y": 1214}, {"x": 387, "y": 1377}, {"x": 330, "y": 1376}, {"x": 332, "y": 1213}]}}, {"description": "America", "boundingPoly": {"vertices": [{"x": 389, "y": 1436}, {"x": 386, "y": 1675}, {"x": 330, "y": 1674}, {"x": 333, "y": 1435}]}}, {"description": "ISBN", "boundingPoly": {"vertices": [{"x": 363, "y": 2501}, {"x": 361, "y": 2558}, {"x": 339, "y": 2557}, {"x": 341, "y": 2500}]}}, {"description": "0", "boundingPoly": {"vertices": [{"x": 360, "y": 2575}, {"x": 360, "y": 2587}, {"x": 338, "y": 2586}, {"x": 338, "y": 2574}]}}, {"description": "14", "boundingPoly": {"vertices": [{"x": 360, "y": 2601}, {"x": 359, "y": 2627}, {"x": 337, "y": 2626}, {"x": 338, "y": 2600}]}}, {"description": "00.2639", "boundingPoly": {"vertices": [{"x": 326, "y": 2501}, {"x": 323, "y": 2604}, {"x": 299, "y": 2603}, {"x": 302, "y": 2500}]}}, {"description": "8", "boundingPoly": {"vertices": [{"x": 323, "y": 2612}, {"x": 323, "y": 2626}, {"x": 299, "y": 2625}, {"x": 299, "y": 2611}]}}, {"description": "Franz", "boundingPoly": {"vertices": [{"x": 214, "y": 1074}, {"x": 208, "y": 1241}, {"x": 147, "y": 1239}, {"x": 153, "y": 1072}]}}, {"description": "Kafka", "boundingPoly": {"vertices": [{"x": 208, "y": 1260}, {"x": 202, "y": 1439}, {"x": 140, "y": 1437}, {"x": 146, "y": 1258}]}}, {"description": "The", "boundingPoly": {"vertices": [{"x": 206, "y": 1476}, {"x": 202, "y": 1599}, {"x": 127, "y": 1596}, {"x": 131, "y": 1474}]}}, {"description": "Trial", "boundingPoly": {"vertices": [{"x": 202, "y": 1616}, {"x": 198, "y": 1751}, {"x": 122, "y": 1748}, {"x": 126, "y": 1613}]}}, {"description": "ISBN", "boundingPoly": {"vertices": [{"x": 156, "y": 2531}, {"x": 154, "y": 2590}, {"x": 131, "y": 2589}, {"x": 133, "y": 2530}]}}, {"description": "0", "boundingPoly": {"vertices": [{"x": 154, "y": 2609}, {"x": 154, "y": 2624}, {"x": 131, "y": 2623}, {"x": 131, "y": 2608}]}}, {"description": "14", "boundingPoly": {"vertices": [{"x": 152, "y": 2643}, {"x": 151, "y": 2667}, {"x": 127, "y": 2666}, {"x": 128, "y": 2642}]}}, {"description": "00.0907", "boundingPoly": {"vertices": [{"x": 124, "y": 2529}, {"x": 121, "y": 2629}, {"x": 97, "y": 2628}, {"x": 100, "y": 2528}]}}, {"description": "8", "boundingPoly": {"vertices": [{"x": 118, "y": 2653}, {"x": 118, "y": 2665}, {"x": 97, "y": 2664}, {"x": 97, "y": 2652}]}}]`

func TestCleanOCR(t *testing.T) {
	assert.Equal(t, "Hi", CleanOCR("Hi ISBN"))
	assert.Equal(t, "Hi", CleanOCR(" # 0123.45 Hi\" -hi"))
}

func TestDimension(t *testing.T) {
	lines, fragments := GetLinesAndFragments(SAMPLE)
	assert.Equal(t, "PMC", lines[0])
	assert.Equal(t, "PMC", fragments[0].Description)
	assert.Equal(t, 333, fragments[0].BoundingPoly.Vertices[0].X)
	assert.Equal(t, 36, MaxDimension(fragments[0].BoundingPoly))
}

func TestPruneSmallText(t *testing.T) {
	lines, fragments := GetLinesAndFragments(SAMPLE)

	// Force some pruning.
	assert.Equal(t, 83, PruneSmallText(lines, fragments, 1))
	assert.Equal(t, 1, PruneSmallText(lines, fragments, PRUNE_SMALL_TEXT))
}

func TestIdentifySpines(t *testing.T) {
	lines, fragments := GetLinesAndFragments(SAMPLE)
	spines := ExtractSpines(lines, fragments)
	log.Printf("Spine %+v", spines[0])
	assert.Equal(t, "PMC", spines[0].Spine)
}