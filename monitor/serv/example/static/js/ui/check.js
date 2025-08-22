const checkValidity = (log) => 
    !!log.edges.evidence_log && !log.edges.is_active

const checkValidities = (logs) => {
    let result = true

    for (const log of logs) {
        result &= checkValidity(log)
    }

    return result
}