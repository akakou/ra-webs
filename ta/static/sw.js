importScripts('./crypto.js')
importScripts('../public_key.js')

const originalFetch = fetch
const TEE_URL = location.origin + '/ra-webs'

self.addEventListener("fetch", async event => {
    event.respondWith(
        customFetch(event.request)
    );
});


const customFetch = async (input, init = {}) => {
    var { url, input } = normalizeURL(input, init)

    if (url.host !== location.host) {
        return originalFetch(input, init)
    }

    console.log('clone')
    const cloned = input.clone()


    const plainReq = {
        url: cloned.url,
        method: cloned.method,
        headers: cloned.headers,
        body: null
    }

    if (input.method === 'POST') {
        const arrayBody = await cloned.arrayBuffer()
        plainReq.body = btoa(String.fromCharCode(...new Uint8Array(arrayBody)))

    }

    const jsonPlainReq = JSON.stringify(plainReq)
    const rawPlainReq = str2ab(jsonPlainReq)

    const { comkey, pubKeyCipher } = await initKey()
    const reqCipher = await encrypt(rawPlainReq, comkey)
    const encodedReqCipher = encodeCipher(reqCipher.cipher, reqCipher.iv, pubKeyCipher)

    const req = new Request(input, {
        method: 'POST',
        body: encodedReqCipher,
        credentials: 'omit',
        headers: {
            'Content-Type': 'application/json'
        },
        url: TEE_URL
    });

    let originalFetching = originalFetch(req, init)
    const originalResp = await originalFetching

    const body = await originalResp.body.text()
    const { cipher, iv } = decodeCipher(body)
    const rawResp = await decrypt(cipher, iv, comkey)


    const resp = new Response(rawResp, {
        status: rawResp.status,
        headers: rawResp.headers,
        body: rawResp.body
    })

    const fetching = async () => resp
    return fetching
}



const normalizeURL = (input, init) => {
    let url;

    if (typeof input === "string") {
        if (input[0] === "/")
            input = `${input.slice(1, input.length)}`;

        if (!input.startsWith("https://")) {
            input = location.host + "/" + input
        }

        url = new URL(input)
        input = new Request(url, init)
    }
    else {
        url = new URL(input.url)
    }

    return { url, input }
}

