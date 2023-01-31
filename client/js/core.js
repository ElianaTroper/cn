
// FUTURE: This would be better off as typescript tbh
// FUTURE: Check for a local IPFS node - if there is one, use that rather than the browser node

const ROOTKEY = '' // TODO: This
const ROOTHASH = 'bafyreiepzzi6vmros67gejuicoflamntixw7mbytetmg3jolcdzoe4cbh4' // TODO: Replaces this with IPNS hash
export const NICENODE = '/dns4/gu1.troper.report/tcp/4002/wss/p2p/12D3KooWES7SHYRm6mNLQAwoycPE1snCZrWkSSY2KibbYhYQWDxf'
const CONFIG = {
	config: {
		Addresses: {
			Swarm: [
				'/dns4/wrtc-star1.par.dwebops.pub/tcp/443/wss/p2p-webrtc-star',
        		'/dns4/wrtc-star2.sjc.dwebops.pub/tcp/443/wss/p2p-webrtc-star'
			]
		}
	},
}
// XXX: The config basically outlines the segments of bootstrapping most prone to censorship - specifically,
//			Preload nodes
//			Delegate nodes
//			Bootstrapping nodes
//
//		Effectively, to obtain some way around censorship of default nodes, we need resistant transports.
//			Once we get working nodes, we should actually be set to make connections

function sleep(ms) {
	// CITATION: https://stackoverflow.com/questions/951021/what-is-the-javascript-version-of-sleep
		// FUTURE: It would be cool to add a citation format for the compiler (ie, these are kinda dead weight if you aren't actually trying to view source)
    return new Promise(resolve => setTimeout(resolve, ms));
}

function btoi(bytes) {
	// CITATION: https://stackoverflow.com/a/69721696
	// XXX: Only works for up to 4 bytes
	let n = 0;
    for (const byte of bytes.values()) {       
            n = (n<<8)|byte;
    }
    return n;
}

class cNode {
	constructor(ipfsNode, rootHash) {
		this.ipfsNode = ipfsNode
		this.rootHash = rootHash
		this.activeCalls = 0
		this.crypto = window.nobleEd25519
	}
	
	async start() {
		this.stabilityMaintinence()
		this.root = await this.get(this.RootHash)
		this.watchRoot()
	}
	
	async get(hash) {
		const cid = Ipfs.CID.parse(hash)
		this.activeCalls += 1
		const res = (await node.dag.get(cid)).value
		this.activeCalls -= 1

		return res
	}
	
	getPrivateKey() {
		// XXX: THIS IS SUPER INSECURE. THESE KEYS ARE STORED IN PLAINTEXT IN LOCALSTORAGE.
		//			DO NOT TRUST THESE KEYS FOR ANYTHING IMPORTANT.
		// For now, this just randomly generates if it's needed
		
		// FUTURE: Something way more robust than this
		if (!this.privateKey) {
			const stored = localStorage.getItem('cNode.privateKey')
			if (stored) {
				const enc = new TextEncoder();
				this.privateKey = enc.encode(stored)
			} else {
				this.privateKey = this.crypto.utils.randomPrivateKey()
				const dec = new TextDecoder()
				localStorage.setItem('cNode.privateKey', dec.decode(this.privateKey))
			}
		}
	}
	
	async sign(message) {
		this.getPrivateKey()
		return await this.crypto.sign(message, this.privateKey)
	}
	
	async verify(publicKey, signature, message) {
		return await this.crypto.verify(signature, message, publicKey)
	}
	
	async watchRoot() {
		// Messages are of the format [index: 4 bytes, hash: the rest]
		this.ipfsNode.subscribe('cn1-test-root',
			async msg => {
				 const hashIndex = btoi(msg.data.slice(64, 68))
				 if (hashIndex > this.root.index) {
				 	const signature = msg.data.slice(0, 64)
				 	const verified = await this.verify(ROOTKEY, signature, msg.data.slice(64))
				 	if (verified) {
				 		newRoot = await this.get(msg.data.slice(4).toString())
				 		if (newRoot.index > this.root.index) {	// Because of async and timing, we double check after the await)
				 			this.root = newRoot
				 		}
				 	} else {
				 		console.log('Bad signature on root pubsub!')
				 	}
				 }
			}
		)
	}
	
	async stabilityMaintinence() {
		// XXX: This greatly speeds up fetching, but it is semi-centralized...
			// The solution would be to get 'dht.findprovs' working properly, but there seem to be many issues there
			// The main drawback is a censor could watch for requests/loads directly to NICENODE
			// FUTUREWORK: Address this (maybe with a censored/paranoid flag?)
		// FUTUREWORK: make activeCalls be a list of hashes, instead of a simple count
			// This will eventually allow for "better" stability maintinence, if dht.findprovs gets sorted...
		while (true) {
			if (this.activeCalls > 0) {
				await this.ipfsNode.swarm.connect(NICENODE)
				await sleep(4000)
			}
			await sleep(1000)
		}
	}
	
	async cat(hash) {
		console.log('cat: started')
		let res = ''
		this.activeCalls += 1
		for await (const chunk of this.ipfsNode.cat(hash)) {
			res += chunk
		}
		this.activeCalls -= 1
		console.log('cat: got content')
		return res
	}
	
	async ls(hash) {
		console.log('ls: started')
		// const resolved = await this.ipfsNode.resolve(hash)
		// console.log(resolved)
		// const provs = this.ipfsNode.dht.findProvs(Ipfs.CID.parse('bafybeicwrjy4peslxwqihmsj7cq2zopswqhv4fvg6dlrf23mfgdoowxygi'))
		// for await (const provider of provs) {
  		//	console.log(provider.id.toString())
		// }
		let res = []
		this.activeCalls += 1
		for await (const file of this.ipfsNode.ls(hash)) {
			res.push(file)
		}
		this.activeCalls -= 1
		console.log('ls: finished')
		return res
	}
	
	async stat(hash) {
		// Stat lists the contents of an IPFS hash
	
		console.log('stat: started')
		this.activeCalls += 1
		const res = await this.ipfsNode.files.stat(hash)
		this.activeCalls -= 1
		console.log('stat: finishes')
		return res
	}
}

console.log('core: creating an ipfs client')
const node = await Ipfs.create(CONFIG)

node.bootstrap.add(NICENODE) // For user who can access this node, we add a known node which contains needed hashes to enable quicker bootstrapping

console.log('core: created a node')

console.log('starting...')

const _NODE = new cNode(node, ROOTHASH)

await _NODE.start()

// FUTURE: Let root update

export const NODE = _NODE
console.log('core: succesfully established a node')
