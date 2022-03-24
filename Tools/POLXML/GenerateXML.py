import xml.etree.cElementTree as ET

root = ET.Element("RegexFormatter")
raw_file = "raw.txt"
out_file = "out.xml"

file = open(raw_file, "r")
while line := file.readline():
	segment = line[9: len(line) - 2] # Remove <!ENTITY... tag
	sep = segment.index(" ")
	find = segment[0:sep]
	replace = segment[sep + 2: -1]
	newpos = ET.SubElement(root, "pos")
	ET.SubElement(newpos, "find").text = find
	ET.SubElement(newpos, "replace").text = replace

e = ET.SubElement(root, "pos")
ET.SubElement(e, "find").text = "&"
ET.SubElement(e, "replace").text = ""

f = ET.SubElement(root, "pos")
ET.SubElement(f, "find").text = ";"
ET.SubElement(f, "replace").text = ""

tree = ET.ElementTree(root)
tree.write(out_file)
