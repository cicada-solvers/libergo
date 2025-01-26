# DWH - Deep Web Hash
## Deep Web Hash - Brute Force Description
The Deep Web Hash (DWH) is a hash that is used in the Liber Primus on page 56 (https://uncovering-cicada.fandom.com/wiki/PAGE_56).  
The hash is a 512-bit hash that is generated from something we don't know at this time.

Since we do not know what the source of the has is, we are going to have to brute for the byte arrays to fine one that 
fits the hash.

### DWH - Brute force toolset
Download the latest release from the releases are of this GitHub.
In the release, there will be 5 files you will need to use.  You do not have to worry about generating arrays.
I am in the process of getting a site ready to give the files to you.

1. generatebytearrays - This is used to generate the byte arrays that are used in the brute force process.
- It has an appsettings.json file that you can use to configure the program.
- num_workers - This is the number of workers that will be used to generate the byte arrays.
- max_permutations_per_line - This is the number of permutations that will be generated per line.
- max_permutations_per_file - This is the number of permutations that will be generated per file.
- max_files_per_zip - This is the number of files that will be generated per zip file.

2. processhashes - This is used to process the hashes that are generated from the byte arrays.
- It has an appsettings.json file that you can use to configure the program.
- num_workers - This is the number of workers that will be used to process the hashes.
- existing_hash - This is the hash that you are looking for. *DO NOT CHANGE THIS*
- Just drop a zip file from the site into the same directory as the processhashes file and run it.

The zip file format is as follows:
package_l(length of arrays)_(zip number)_\_of_(total number of zips).zip

Once you have downloaded a .zip file from the site (TBD), then you will need to run ./processhashes.
It will find the zip files in the directory and then start hashing the byte array into SHA-512, Whirlpool, and Blake2b-512.
It will remove the range from the permutation text file that is in the zip file.  Then it will remove the zip file.
If you stop the process, it will resume on the last permutation file you were processing.  That way, you do not have to 
tie up your machine for weeks on end.

If you have others you would like to try THAT WERE OUT AT THE TIME OF THE PUZZLE, then feel free to hit me up on the
[Cicada Solvers Discord](https://discord.com/invite/5qznJtjw?utm_source=Discord%20Widget&utm_medium=Connect).
I lurk there quite often.