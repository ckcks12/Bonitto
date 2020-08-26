import React, {FC, useCallback, useEffect, useLayoutEffect, useState}      from "react";
import {Button, Col, Divider, List, message, Radio, Row, Spin, Typography} from "antd";
import {Boilerplate, Problem, Result}                                      from "../lib/type";
import {useParams}                                                         from "react-router-dom";
import {ExclamationOutlined}                                               from "@ant-design/icons";
import {APIGetProblem, APIGetResult, APISubmit}                            from "../lib/api";
import AceEditor                                                           from "react-ace";
import * as faker                                                          from "faker";
import "ace-builds/src-noconflict/mode-javascript";
import "ace-builds/src-noconflict/mode-golang";
import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/mode-c_cpp";
import "ace-builds/src-noconflict/theme-terminal";
import "ace-builds/src-noconflict/ext-language_tools";

function bindWebSocket(url: string, cb: Function) {
    const ws = new WebSocket(url);
    let flag = false;
    ws.onclose = () => {
        if (flag) return;
        flag = true;
        message.warn("Disconnected from server, reconnecting...");
        setTimeout(() => bindWebSocket(url, cb), 3000);
    };
    ws.onerror = () => {
        //@ts-ignore
        ws.onclose?.();
    };
    cb(ws);
}

function getOrGenerateId(): string {
    let id = new URLSearchParams(window.location.search).get("id");
    if (id) return id;
    return faker.name.firstName() + faker.name.lastName();
}

export const ProblemPage: FC = () => {
    const {no} = useParams();
    const [problem, setProblem] = useState<Problem>();
    const [code, setCode] = useState("");
    const [boilerplate, setBoilerplate] = useState<Boilerplate>({code: "", lang: "javascript"});
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [submittingMsg, setSubmittingMsg] = useState("");
    const [id, setId] = useState(getOrGenerateId());
    const [ws, setWs] = useState<WebSocket>();
    const [result, setResult] = useState<Result[]>([]);

    useLayoutEffect(() => {
        (async () => {
            try {
                setProblem(await APIGetProblem(no));
                setResult(await APIGetResult(id, no));
            } catch (e) {
                message.error(e);
            }
        })();
        bindWebSocket(`ws://${window.location.hostname}/api/ws/${id}`, setWs);
    }, []);

    useEffect(() => {
        if (!problem) return;
        setBoilerplate(problem.boilerplate[0]);
    }, [problem]);

    useEffect(() => {
        setCode(boilerplate.code);
    }, [boilerplate]);

    useEffect(() => {
        if (!ws) return;
        ws.onopen = (e) => {
            message.info(`Welcome, ${id}`);
            setInterval(() => ws?.send(""), 500);
        };
        ws.onmessage = (e) => {
            setSubmittingMsg(e.data);
            if (e.data.indexOf("😓") > -1 || e.data.indexOf("🤩") > -1) {
                setTimeout(afterSubmit, 3000);
            }
        };
    }, [ws]);

    const afterSubmit = useCallback(async () => {
        setIsSubmitting(false);
        setResult(await APIGetResult(id, no));
        message.success(`You've got new result, ${id}`);
    }, [id, no]);

    const submit = useCallback(async () => {
        setSubmittingMsg("Sending...");
        setIsSubmitting(true);
        await APISubmit(no, id, code, boilerplate.lang);
    }, [id, code, boilerplate]);

    const handleSetLang = useCallback((e) => {
        if (!problem) return;
        const lang = e.target.value;
        const b = problem.boilerplate.find((b) => b.lang === lang);
        b && setBoilerplate(b);
    }, [problem]);

    const renderResult = useCallback((l: Result) => (
        <List.Item>
            <List.Item.Meta
                avatar={<ExclamationOutlined/>}
                title={l.result}
                description={l.date.toLocaleString()}
            />
        </List.Item>
    ), []);

    return (
        <Row>
            <Col xs={24} sm={24} md={22} lg={18} xl={16} style={{margin: "auto"}}>
                <Typography.Title>{problem?.no}. {problem?.title}</Typography.Title>
                <Divider/>
                <Typography.Paragraph>
                    <pre>{problem?.content}</pre>
                </Typography.Paragraph>
                <Divider>{id}'s Code</Divider>
                <Radio.Group value={boilerplate.lang} onChange={handleSetLang}
                             style={{display: "block", textAlign: "center"}}>
                    <Radio.Button value="javascript">javascript</Radio.Button>
                    <Radio.Button value="java" disabled>java</Radio.Button>
                    <Radio.Button value="golang">golang</Radio.Button>
                    <Radio.Button value="cpp" disabled>cpp</Radio.Button>
                </Radio.Group>
                <br/>
                <Spin tip={submittingMsg} spinning={isSubmitting}>
                    <AceEditor
                        mode={boilerplate.lang}
                        theme="terminal"
                        name="code"
                        fontSize={16}
                        showPrintMargin={true}
                        showGutter={true}
                        highlightActiveLine={true}
                        setOptions={{
                            enableBasicAutocompletion: true,
                            enableLiveAutocompletion: true,
                            showLineNumbers: true,
                            useWorker: false,
                            displayIndentGuides: true,
                        }}
                        style={{width: "100%"}}
                        value={code}
                        onChange={setCode}
                    />
                    <br/>
                    <Button onClick={submit} style={{display: "block", margin: "auto"}}>Submit</Button>
                </Spin>
                <Divider>{id}'s permalink</Divider>
                <Typography.Text
                    copyable={true}>{`${window.location.protocol}${window.location.host}${window.location.pathname}?id=${id}`}</Typography.Text>
                <Divider>{id}'s Last Submit</Divider>
                <List size="small" bordered dataSource={result}
                      renderItem={renderResult}
                      style={{lineBreak: "anywhere"}}
                />
            </Col>
        </Row>
    );
};
