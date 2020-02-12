#!/bin/sh
##
## batch.sh
## 
## Made by 
## Login   <clinares@atlas>
## 
## Started on  Sun Dec 16 21:20:46 2018 
## Last update Sun Dec 16 21:20:46 2018 
##

for ifile in `ls *.tex`
do
    pdflatex $ifile
    pdflatex $ifile
    pdflatex $ifile
done

rm *.aux *.log
