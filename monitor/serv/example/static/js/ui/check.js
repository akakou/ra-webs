const checkValidity = (Log) => {
    let result = true

    for (const log of Log) {
        result &= !log.violation
    }

    return result
}