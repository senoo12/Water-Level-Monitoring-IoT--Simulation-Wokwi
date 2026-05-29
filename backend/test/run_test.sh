#!/bin/bash

echo "===== START TESTING ====="
echo ""

go test ./test -run Test_BeforeAuthentication -v

echo ""
echo "=============================="
echo ""

go test ./test -run Test_AfterAuthentication_WithoutToken -v

echo ""
echo "=============================="
echo ""

go test ./test -run Test_AfterAuthentication_WithValidToken -v

echo ""
echo "=============================="
echo ""

go test ./test -run Test_AfterAuthentication_InvalidToken -v

echo ""
echo "=============================="
echo ""

go test ./test -run Test_BruteForce_InvalidToken -v

echo ""
echo "===== TESTING FINISHED ====="