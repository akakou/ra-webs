const TableCompornent = ({ logs }) => {
    const rows = logs.map((server, index) =>
        <tr key={index}>
            <td>
                {index + 1}
            </td>
            <td>{server.domain}</td>
            <td>{server.edges.code.repository}</td>
            <td>{server.edges.code.commit_id}</td>
            <td>{server.edges.code.unique_id}</td>
            <td>{!!server.has_activated}</td>
            <td><a
                href={`${server.edges.code.repository}/tree/${server.edges.code.commit_id}`}
            >Go to Repo</a></td>
        </tr>
    );

    return (
        <table>
            <thead>
                <tr>
                    <th>Index</th>
                    <th>Domain</th>
                    <th>Repository</th>
                    <th>Commit_ID</th>
                    <th>Unique ID</th>
                    <th>Activated</th>
                    <th>Repository Link</th>
                </tr>
            </thead>
            <tbody>
                {rows}
            </tbody>
        </table>
    )
}