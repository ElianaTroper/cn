const HASH = 'bafybeiaysi4s6lnjev27ln5icwm6tueaw2vdykrtjkwiphwekaywqhcjze'
const PATH = '/58wiki'
const node = await Ipfs.create()                               
console.log('core: created a node')
const args = { path: PATH }
//const args = { path: '/' }
/*
console.log(args)
console.log(node.dag.get(await Ipfs.CID.parse(res)))
const hash = await node.dag.resolve(Ipfs.CID.parse(res), args)
console.log(hash)
const hash2 = await node.dag.get(hash.cid)
console.log(hash2)
*/

/*
for await (const name of node.name.resolve('/ipfs/bafybeiaysi4s6lnjev27ln5icwm6tueaw2vdykrtjkwiphwekaywqhcjze/wiki/Cinematography')) {
  console.log(name)
  // /ipfs/QmQrX8hka2BtNHa8N8arAq16TCVx5qHcb46c5yPewRycLm
}
*/

let curr = '/ipfs/bafybeiaysi4s6lnjev27ln5icwm6tueaw2vdykrtjkwiphwekaywqhcjze/wiki'
console.log(curr)
curr = await node.resolve(curr)
console.log(curr)

// const t = await node.resolve('/ipfs/QmPzZpDqsXeeLt4vEB7TuVs622jp5ECHNeKGDxoMxDDDPW/wiki/Books')

console.log(t)

const x = await node.get('/ipfs/bafybeiaysi4s6lnjev27ln5icwm6tueaw2vdykrtjkwiphwekaywqhcjze/wiki/Book')


for await (const data of node.get('/ipfs/QmPzZpDqsXeeLt4vEB7TuVs622jp5ECHNeKGDxoMxDDDPW/wiki/Books')) {
	console.log(data)
	res += data
}


/*
res = node.dag.get('/ipfs/bafybeiaysi4s6lnjev27ln5icwm6tueaw2vdykrtjkwiphwekaywqhcjze/wiki/Cinematography')
*/

console.log(res)
