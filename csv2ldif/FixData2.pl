# -
# -  FixData.pl
# -
# -	This program scans CSV export data and maps it to an update data set for ldap
# -
# -     You will be prompted for the file to convert as well as the attribute change
# -     type.  The CSV should be in a format where the first row has column headings
# -     and the first column should be the DN.  The output file will be the name of
# -     the input file with the LDIF extension added.
# -
# -
# -  Example:  perl FixData.pl
# -

use FileHandle;

my($csvFile, $ldifType, $inputfile, %domains, $outfile);
my(@linearray, @dnarray, @cnarray, $linecount);
my(%controlhash, $templatefile, %fieldmap, $logfile);


$csvFile = &promptUser("Name of file to be converted from CSV to LDIF 'Add' File");
$ldifType = &promptUser("Will these attributes be in replace or add mode?", "replace");
$inputfile = new FileHandle($csvFile) ||
	die "Can't open the import file!\n";
$outfile = new FileHandle(">$csvFile.ldif") ||
	die "Can't create or open the output file!\n";
$logfile = new FileHandle(">FixData.log") ||
	die "Can't create or open the log file!\n";

# - Read the first line of the data file an build the control hash from it.


print $outfile "\n";
print $outfile "version: 1\n\n";

$linecount = 0;
while(defined($line = <$inputfile>)) {
	$linecount++;
	chomp($line);
	if($linecount == 1) {
		@linearray = split(m',', $line);
		$fieldcount = 0;

		foreach $field (split(m',', $line)) {
			$fieldmap{$fieldcount} = $field;
			$fieldcount++;
		}
		next;
	}

	@linearray = split(m'"', $line);
	@fieldarray = split(m',', $linearray[2]);

	for($k = 0; $k < $fieldcount; $k++) {
		if($k == 0) {
			print $outfile "$fieldmap{$k}: ";
			print $outfile "$linearray[1]\n";
			print $outfile "changetype: modify\n";
		} else {
			if($k > 1) {
				print $outfile "-\n";
			}
			print $outfile "$ldifType: ";
			print $outfile "$fieldmap{$k}\n";
			print $outfile "$fieldmap{$k}: $fieldarray[$k]\n";
		}
	}
	print $outfile "\n";
}

sub promptUser {

   #-------------------------------------------------------------------#
   #  two possible input arguments - $promptString, and $defaultValue  #
   #  make the input arguments local variables.                        #
   #-------------------------------------------------------------------#

   local($promptString,$defaultValue) = @_;

   #-------------------------------------------------------------------#
   #  if there is a default value, use the first print statement; if   #
   #  no default is provided, print the second string.                 #
   #-------------------------------------------------------------------#

   if ($defaultValue) {
      print $promptString, "[", $defaultValue, "]: ";
   } else {
      print $promptString, ": ";
   }

   $| = 1;               # force a flush after our print
   $_ = <STDIN>;         # get the input from STDIN (presumably the keyboard)


   #------------------------------------------------------------------#
   # remove the newline character from the end of the input the user  #
   # gave us.                                                         #
   #------------------------------------------------------------------#

   chomp;

   #-----------------------------------------------------------------#
   #  if we had a $default value, and the user gave us input, then   #
   #  return the input; if we had a default, and they gave us no     #
   #  no input, return the $defaultValue.                            #
   #                                                                 # 
   #  if we did not have a default value, then just return whatever  #
   #  the user gave us.  if they just hit the <enter> key,           #
   #  the calling routine will have to deal with that.               #
   #-----------------------------------------------------------------#

   if ("$defaultValue") {
      return $_ ? $_ : $defaultValue;    # return $_ if it has a value
   } else {
      return $_;
   }
}
