## Useful Queries

// GetSystem
```cypher
MATCH p=shortestPath((low:Principal)-[*..5]->(hi:Principal))
WHERE any(sp in ['system', 'trusted'] where hi.name contains(sp)) and
 none(sp in ['system', 'trusted'] where low.name contains(sp)) 
RETURN p
```

// Lateral Movement
```cypher
MATCH p=allshortestPaths((low:Principal)-[*]->(hi:Principal))
WHERE none(sp in ['system', 'trusted', 'creator', 'any', 'author'] where hi.name contains(sp)) AND
 none(sp in ['system', 'trusted', 'creator', 'any', 'author'] where low.name contains(sp)) AND
 none(r in relationships(p) where type(r) ='IMPORTS')
RETURN p
```

// Dll Hijacks
```cypher
MATCH (p:Principal {name: 'deathstar/alex'})-[*..2]->(pe:INode)<-[:IMPORTED_BY]-(dep:Dep)
where not pe:Dll
and none(known in ["wow64cpu.dll", "wowarmhw.dll", "xtajit.dll", "advapi32.dll", "clbcatq.dll", "combase.dll", "comdlg32.dll", "coml2.dll", "difxapi.dll", "gdi32.dll", "gdiplus.dll", "imagehlp.dll", "imm32.dll", "kernel32.dll", "msctf.dll", "msvcrt.dll", "normaliz.dll", "nsi.dll", "ole32.dll", "oleaut32.dll", "psapi.dll", "rpcrt4.dll", "sechost.dll", "setupapi.dll", "shcore.dll", "shell32.dll", "shlwapi.dll", "user32.dll", "wldap32.dll", "wow64.dll", "wow64win.dll", "ws2_32.dll", "ntdll.dll"] where dep.name = known)
and not (:Directory {path: pe.parent})-[:CONTAINS]->(:Dll {path: pe.parent + '/' + dep.name})
return apoc.map.fromLists(["exe", "imports", "path"],[pe.name, collect(distinct dep.name), pe.parent]) AS hijacks
```

// Show who imports `wer.dll`
```cypher
match (n)-[i:IMPORTS]->(d:DLL {name: "wer.dll"}) return n.path, i.fn, d.path
```

// Find EXEs with dump-related imports in AppData:
```cypher
match (e:EXE)-[fn:IMPORTS]->(d:DLL) 
where fn.fn contains "Dump" 
 and e.path contains "AppData" 
return e.path,fn,d.name
```


// Get EXEs that potentially start RPC servers
```cypher
match a=(e:EXE)-[r:IMPORTS]->(d:DLL {name: "RPCRT4.dll"}) 
where r.fn contains("Binding") 
return e.name,e.path,collect(r.fn) as importedFns
```


// Get RPC server PEs
```cypher
match a=(e)-[r:IMPORTS]->(d) 
where r.fn contains("RpcServerListen") 
 and not e.path contains("System32")
 //and not r.fn contains("auth") 
return e.name,e.path
```


// Get RPC Client PEs:
```cypher
match a=(e)-[r:IMPORTS]->(d) 
where r.fn contains("RpcStringBindingCompose") 
 and not e.path contains("System32")
return e.name,e.path
```

