#
# makefile
#  
# Made by Carlos Linares
# Login   <clinareskorf.plg.inf.uc3m.es>
#
# Started on  <Wed Jan 30 12:33:00 2008 Carlos Linares>
# Last update <martes, 23 febrero 2016 20:37:47 Carlos Linares Lopez (clinares)> 
# -----------------------------------------------------------------------------

# -----------------------------------------------------------------------------
# This makefile considers the most common options for converting (standard) .tex
# files into .pdf files. The procedure is:
# 
# 1. Rename the macro $NAME to the final filename of the pdf file
# 2. Usually, leave the $SRC macro with the "main.tex" file but change it if
#    necessary
# 3. If letter size is requested, add "-t letter" to the suffix rule ".dvi.ps"
# 
# -----------------------------------------------------------------------------

NAME 	= roberto.pdf
TARBALL = roberto.tar.gz
SRC	= main.tex
OBJ     = main.pdf

# Macro Definitions 
# ------------------------------------------
LATEX 	= ./pytex.py
DVIPS   = dvips
PSPDF   = ps2pdf
RM 	= /bin/rm -f

.SUFFIXES: .tex .pdf

.tex.pdf:
	$(LATEX) $<

##############################
# Basic Compile Instructions #
##############################

all:	$(NAME)

$(NAME): $(OBJ)
	@mv $(OBJ) $(NAME)

: $(SRC)

tarball:
	@$(TAR) zcvf $(TARBALL)

clean:
	@$(RM) *.aux *.log *.nav *.out *.snm *.toc *.dvi *.bbl *.blg *.spl *.vrb *~

delete: clean
	@$(RM) $(NAME) $(OBJ) *.ps



/* Local Variables: */
/* mode:makefile */
/* fill-column:80 */
/* End: */
