UNAME := $(shell uname)
default: buildpb

lintf:
	prototool lint .
	prototool format -w

genproto:
	prototool all

buildpb: lintf genproto
