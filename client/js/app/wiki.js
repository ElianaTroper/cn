
function updateLinks(NODE, page, urlParams) {
	for (const link of page.links) {
		// console.log(link)
		const original = link.getAttribute('href')
		// console.log(original)
		// TODO: Remove first .., format properly
		// link.href = "javascript:void(0);"
	}
}

export async function runWiki(NODE, urlParams) {
	// TODO: Clean path of #. ?
	//
	console.log('runWiki: started')

	const path = urlParams.get('path')
	console.log(path)
	let mod = path
	if (!path) {
		mod = 'index.html'
	}
	const fetchPath = HASH + '/wiki/' + mod
	console.log('runWiki: getting ' + fetchPath)

	const page = await NODE.doCat(fetchPath)
	console.log('runWiki: succesfully retrieved the file')
	
	const newPage = new DOMParser().parseFromString(page, 'text/html')
	
	updateLinks(NODE, newPage, urlParams)

	document.getElementsByTagName("html")[0].innerHTML = newPage.documentElement.innerHTML
	
	// TODO: Dynamically insert javascript code and other IPFS resources
	// 	This may involve changing up links, etc based on relative paths	

	// TODO: Replace the body/header with the needed items
	// 		Good news - the script doesn't unload!<F3>
	
	// Eventually, we could replace links with javascript functions
	// 	Basically, we want to maintain the state of the app so
	// 	content can be served from the same node, to keep stuff
	// 	running
}
