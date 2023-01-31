
import { runIndex } from '../main.js';

const PAGES = '/ipfs/bafybeidsjjvzeocm4fihzty7fgjdxhbkqcwh6sqholx7al6apz3alca2zi' // TODO: Get this dynamically

function getBookPath(hash, item){
	let res = '/'
	if (item.length == 1) {
		// Special case: the first 9 items
		res += '/0'
	} else {
		res += item.split('').slice(0, -1).join('/')
	}
	return hash + res + '/' + item
}

async function runBaseIndex(node, params) {

	const page = params.get('page') || '1'
	const pageInfo = JSON.parse(await node.cat(PAGES + '/' + page + '.json'))
	let toIndex = {}
	for (const [_key, value] of Object.entries(pageInfo.items)) {
		toIndex[value['Text#']] = value.Title
	}
	runIndex(node, params, toIndex, 'item')
}

async function runItemIndex(node, params, path, mod='') {
	console.log(path)
	const list = await node.ls(path)
	console.log(list)
	let toIndex = {}
	list.forEach((element) => {
		toIndex[element.name] = element.name;
	})
	runIndex(node, params, toIndex, 'path', mod)
}

async function loadTxt(node, params, cid) {
	const txt = await node.cat(cid)
	document.body.innerHTML = `
		<div style="text-align: center;">
			<pre style="word-wrap: break-word; white-space: pre-wrap; display: inline-block; text-align: left; tab-size: 2; width: 80ch;">
				${txt}
			</pre>
		</div>
		`
}

function runContentPage(node, params, cid) {
	const type = cid.split('.').pop();
	
	switch (type) {
		case 'txt':
			loadTxt(node, params, cid)
			break;
		case 'htm':
			// TODO: Finish
		default:
			alert(type + ' is not (yet) a supported file type!')
	}
}

export async function runGutenberg(node, params) {
	// TODO: Clean path of #. ?

	console.log('runGutenberg: started')

	const item = params.get('item')
	if (!item) {
		runBaseIndex(node, params)
	} else {
		const path = params.get('path')
		if (!path) {
			runItemIndex(node, params, getBookPath(node.root.app.gutenberg.hash, item))
		} else {
			const cid = getBookPath(node.root.app.gutenberg.hash, item) + '/' + path
			console.log(cid)
			const info = await node.stat(cid)
			console.log(info)
			if (info.type === 'file') {
				runContentPage(node, params, cid)
			} else if (info.type === 'directory') {
				// TODO: Need to append to path
				runItemIndex(node, params, cid, path + '/')
			} else {
				console.log(info)
				throw 'Bad item!'
			}
		}
	}
	
	/* OLD WIKI STUFF
	
	const fetchPath = HASH + '/wiki/' + mod
	console.log('runWiki: getting ' + fetchPath)

	const page = await node.doCat(fetchPath)
	console.log('runWiki: succesfully retrieved the file')
	
	const newPage = new DOMParser().parseFromString(page, 'text/html')
	
	updateLinks(node, newPage, params)

	document.getElementsByTagName("html")[0].innerHTML = newPage.documentElement.innerHTML
	
	// TODO: Dynamically insert javascript code and other IPFS resources
	// 	This may involve changing up links, etc based on relative paths	

	// TODO: Replace the body/header with the needed items
	// 		Good news - the script doesn't unload!<F3>
	
	// Eventually, we could replace links with javascript functions
	// 	Basically, we want to maintain the state of the app so
	// 	content can be served from the same node, to keep stuff
	// 	running
	
	*/
}
