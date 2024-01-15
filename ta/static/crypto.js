const str2ab = (str) => {
    const buf = new ArrayBuffer(str.length);
    const bufView = new Uint8Array(buf);
    for (let i = 0, strLen = str.length; i < strLen; i++) {
        bufView[i] = str.charCodeAt(i);
    }
    return buf;
}

const IV_SIZE = 96 / 8

const initKey = async () => {
    const comkey = await crypto.subtle.generateKey(
        {
            name: "AES-GCM",
            length: 256,
        },
        true,
        ["encrypt", "decrypt"],
    );

    const exportedComkey = await crypto.subtle.exportKey("raw", comkey);
    const rawComkey = new Uint8Array(exportedComkey);

    const base64Pubkey = atob(PUBLIC_KEY);
    const derPubkey = str2ab(base64Pubkey);

    const pubKey = await crypto.subtle.importKey(
        "spki",
        derPubkey,
        {
            name: "RSA-OAEP",
            hash: "SHA-256",
        },
        true,
        ["encrypt"],
    );

    const pubKeyCipher = await crypto.subtle.encrypt(
        {
            name: "RSA-OAEP",
        },
        pubKey,
        rawComkey,
    );

    return {
        comkey,
        pubKeyCipher
    }
}


const encrypt = async (plain, comkey) => {
    const iv = crypto.getRandomValues(new Uint8Array(IV_SIZE));

    const cipher = await crypto.subtle.encrypt(
        {
            name: "AES-GCM",
            iv: iv.buffer,
        },
        comkey,
        plain,
    );

    return {
        cipher,
        iv,
    }
}


const decrypt = async (cipher, comkey) => {
    const iv = crypto.getRandomValues(new Uint8Array(IV_SIZE));

    const plain = await crypto.subtle.decrypt(
        {
            name: "AES-GCM",
            iv: iv.buffer,
        },
        comkey,
        cipher,
    );

    return plain
}


const encodeCipher = (cipher, iv, pubKeyCipher) => {
    const base64Cipher = btoa(String.fromCharCode(...new Uint8Array(cipher)))
    const base64IV = btoa(String.fromCharCode(...new Uint8Array(iv.buffer)))
    const base64PubKeyCipher = btoa(String.fromCharCode(...new Uint8Array(pubKeyCipher)))

    return JSON.stringify({
        iv: base64IV,
        text: base64Cipher,
        key: base64PubKeyCipher
    })
}

const decodeCipher = ({ iv, cipher }) => {
    // const { cipher, iv } = JSON.parse(allCipher)
    const decodedCipher = str2ab(atob(cipher))
    const decodedIV = str2ab(atob(iv))

    return {
        text: decodedCipher,
        iv: decodedIV
    }
}