import React, {FC, useLayoutEffect, useState} from "react";
import {Col, Row, Table, Typography}          from "antd";
import {Problem}                              from "../lib/type";
import {Link}                                 from "react-router-dom";
import {APIGetProblems}                       from "../lib/api";

export const ProblemListPage: FC = () => {
    const [problems, setProblems] = useState<Problem[]>([]);

    useLayoutEffect(() => {
        (async () => {
            setProblems(await APIGetProblems());
        })();
    }, []);

    return (
        <Row>
            <Col xs={24} sm={24} md={22} lg={18} xl={16} style={{margin: "auto"}}>
                <Typography.Title>Problems</Typography.Title>
                <Table
                    size="small"
                    bordered
                    dataSource={problems}
                >
                    <Table.Column title="No" key="no" dataIndex="no" width="10%"/>
                    <Table.Column title="Title" key="title" dataIndex="title" render={(_, p: Problem) => (
                        <Link to={`/problems/${p.no}`}>{p.title}</Link>
                    )}/>
                </Table>
            </Col>
        </Row>
    );
};
