const checkOneValidity = x => !x.edges.violation

const checkAllValidity = (data) => {
    if (!data.length) return false

    const reduced = data.reduce((accumulator, current) => accumulator && checkOneValidity(current), true)
    const last_activated = data[data.length - 1].has_activated
    return last_activated && reduced
}


const fetchAndValidate = async (url) => {
    const resp = await axios.get(`/api/ta/${url}`)
    console.log(resp)

    const data = resp.data

    const v = data.length > 0 && checkAllValidity(data)
    return {v, data}
}