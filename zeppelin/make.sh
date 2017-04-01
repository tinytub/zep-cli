mv chash.cc.bak chash.cc
rm chash.o lib/libchash.a
gcc-4.9 -c chash.cc -std=c++11 
#gcc -c chash.cc -std=c++11
/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/ar rv libchash.a chash.o
mv libchash.a lib
mv chash.cc chash.cc.bak
