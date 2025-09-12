# Group Membership and Privilege Escalation

## Overview

The `lpegopher` tool now properly supports `MEMBER_OF` relationships between users and groups, enabling detection of privilege escalation paths through group membership.

## How it Works

### Data Structure

- **Users and Groups**: Both represented as `Principal` nodes in Neo4j
- **Group Property**: User principals have a `group` property that references their group membership
- **MEMBER_OF Relationship**: Created between user and group principals when `user.group = group.name`

### Relationship Creation

The system creates `MEMBER_OF` relationships during the processing phase:

```cypher
MATCH (group:Principal),(user:Principal) 
WHERE user.group = group.name 
RETURN user, group
```

This creates: `(user)-[:MEMBER_OF]->(group)`

### Privilege Escalation Detection

The membership relationships integrate with existing privilege escalation queries. For example, the "GetSystem" query:

```cypher
MATCH p=shortestPath((low:Principal)-[*..5]->(hi:Principal))
WHERE any(sp in ['system', 'trusted'] where hi.name contains(sp)) and
 none(sp in ['system', 'trusted'] where low.name contains(sp))
RETURN p
```

This query will now find paths like:
1. `(user)-[:MEMBER_OF]->(group)-[:OWNS]->(system_resource)`
2. `(user)-[:MEMBER_OF]->(group)-[:WRITE_DACL]->(privileged_file)`

### CSV Data Format

For group membership to work, the `principals.csv` file should contain:

```csv
nid,name,group
user_id_1,alice,Developers
user_id_2,bob,Administrators  
group_id_1,Developers,
group_id_2,Administrators,
```

- **Users**: Have both `name` and `group` properties populated
- **Groups**: Have `name` populated but `group` can be empty

## Bug Fix

Fixed a critical bug in the processing pipeline where membership relationships might not be created if ACL processing failed. The controller now properly checks for errors after each processing step, ensuring that:

1. If ACL processing fails, the error is returned immediately
2. If membership processing fails, the error is returned immediately  
3. Both operations run to completion if no errors occur

## Testing

The implementation includes comprehensive tests that verify:

- MEMBER_OF relationship template correctness
- Proper relationship direction (user → group)
- Integration with privilege escalation scenarios
- Error handling in the processing pipeline

Run tests with:
```bash
go test ./node -v -run="TestMembership|TestPrivilege"
```