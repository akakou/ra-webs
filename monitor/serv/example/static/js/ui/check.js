const checkValidity = (log) => 
    !!log.edges.at_log && !log.edges.violation

const checkValidities = (logs) => {
    let result = true

    for (const log of logs) {
        result &= checkValidity(log)
    }

    return result
}