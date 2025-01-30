# What do?
rewrite uids, gids, usernames, and groupnames inside of tar streams 

# Usage

## Create a new tarball
```cat tarball.tar | chowntar [options] > new-tarball.tar```

## Extract to the filesystem
```cat tarball.tar | chowntar [options] | tar -vx```

## Multiple rewrites
```cat tarball.tar | chowntar [options] | chowntar [options] | chowntar [options]```

# Options
```
Usage of ./chowntar:
  -from-gid int
    	Rewrite GID from GID
  -from-group string
    	Rewrite groupname FROM groupname
  -from-uid int
    	Rewrite UID from UID
  -from-user string
    	Rewrite username FROM username
  -to-gid int
    	Rewrite this GID to GID
  -to-group string
    	Rewrite groupname FROM groupname
  -to-uid int
    	Rewrite UID to  UID
  -to-user string
    	Rewrite username FROM username
  -verbose
    	Be verbose about changes
```
