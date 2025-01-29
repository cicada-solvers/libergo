# DWH - Deep Web Hash
## Deep Web Hash - Brute Force Description
The Deep Web Hash (DWH) is a hash that is used in the Liber Primus on page 56 (https://uncovering-cicada.fandom.com/wiki/PAGE_56).  
The hash is a 512-bit hash that is generated from something we don't know at this time.

Since we do not know what the source of the hash is, we are going to have to brute for the byte arrays to fine one that 
fits the hash.  Once we have the byte array, we can figure it out.

## DWH - Brute force toolset
Download the latest release from the releases on this GitHub repo.
In the release, there will be two programs you will need to use.  Both of these will be under cmd.

1. generatebytearrays - This is used to generate the byte arrays that are used in the brute force process.
- It has an appsettings.json file that you can use to configure the program.
- num_workers - This is the number of workers that will be used to generate the byte arrays.
- max_permutations_per_line - This is the number of permutations that will be generated per range (line).
- max_ranges_per_segment - This is the number of permutations ranges that will be generated per segment.
- max_segments_per_package - This is the number of segments that will be generated per package.


2. processhashes - This is used to process the hashes that are generated from the byte arrays.
- It has an appsettings.json file that you can use to configure the program.
- num_workers - This is the number of workers that will be used to process the hashes.  You will want to adjust this for your machine!
- existing_hash - This is the hash that you are looking for. *DO NOT CHANGE THIS*

The ranges in the permutations are 500 million per line. There are 62500 ranges per package.  

There is no expectation that the average PC can chew through one in a day unless you have a thread-ripper or something similar.
It should take several days/weeks to process one package (again, depending on the machine).

I track those willing to progress in a sheet at (https://docs.google.com/spreadsheets/d/1yy1Eg3QLrCloZ9N9Dhfh8FLcx_Po56SCOo49Y-UC068/edit?usp=sharing).
If you are willing to help track people's progress, please let me know.  I will add you as an editor.

### DWH - Brute force process
*Note:You will need to modify the number of worker threads in the appsettings.json file to match your machine's capabilities and CPU utilization desires.*

*Warning: Do not change the folder layout.  It will cause the programs to not work correctly*

1. Check the sheet to see what packages are not currently being processed.
2. On the command line, navigate to the cmd/generatebytearrays folder.
3. Run the generatebytearrays program to generate the byte arrays.
- You will be prompted for the array length and the packages to create.  You can create multiple packages at once. 
4. Once the file has been created, it will save to a parent directory in a database file.
5. Go to the cmd/processhashes folder on the command line.
6. Run the processhashes program to process the hashes.

The hasher will mark the processed line in the database once it has been processed.  This will allow you some degree of resuming the process if you need to stop it for some reason.

### DWH - Brute force results
- If you find the hash, please post it on the forum post.  This will allow others to know that the hash has been found.  You will see it in the found_hashes.txt file.
- If you no longer want to participate in the brute force process, please post on the forum so others can pick up the processing of the file you allocated.

## DWH - Hashes Being Tested
- SHA-512
- Blake2b-512
- Whirlpool

If you have others you would like to try THAT WERE OUT AT THE TIME OF THE PUZZLE, then feel free to hit me up on the
[Cicada Solvers Discord](https://discord.com/invite/5qznJtjw?utm_source=Discord%20Widget&utm_medium=Connect).
I lurk there quite often.

## FAQ
- Q: Why are you doing this?
- A: I am doing this because I want to know what the hash is.  I am also doing this because I want to help the community solve the puzzle.
---
- Q: Is this required to be a part of this community?
- A: No, this is not required.  This is a voluntary effort to help solve the puzzle.
---
- Q: How can I help?
- A: You can help by running the programs on your machine.  You can also help by tracking the progress of the packages in the Google Sheet.
---
- Q: Will I be compensated for my efforts?
- A: No, you will not be compensated for your efforts.  This is a voluntary effort to help solve the puzzle.
---
- Q: What happens if I find the hash?
- A: If you find the hash, please post it on the forum post.  This will allow others to know that the hash has been found.  You will see it in the found_hashes.txt file.