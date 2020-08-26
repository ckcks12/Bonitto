import React, {FC}            from "react";
import {Col, Row, Typography} from "antd";
import styled                 from "styled-components";

const Question = styled.div`
  margin-top: 50px;
  font-size: 1.25rem;
  font-style: italic;
  &::before {
    margin-right: 15px;
    content: "🎤";
  }
`;

const Answer = styled.div`
  font-size: 1.5rem;
  font-weight: bold;
  color: hsla(0, 0%, 100%, .85);
  &::before {
    margin-right: 15px;
    content: "📣";
  }
`;

export const MainPage: FC = () => {
    return (
        <Row>
            <Col xs={24} sm={24} md={22} lg={18} xl={16} style={{margin: "auto"}}>
                <Typography.Title>🍄 Bonitto</Typography.Title>

                <Question>What is Bonitto?</Question>
                <Answer>Bonitto is a website for practicing coding skills.</Answer>

                <Question>What's the difference from others?</Question>
                <Answer>
                    Yes, there are already similar websites.<br/>
                    Bonitto is, however, providing <u><i>Real World Problems</i></u>.
                </Answer>

                <Question><u><i>Real World Problems</i></u>? What is it?</Question>
                <Answer>
                    <u><i>Real World Problems</i></u> are:
                    <ul>
                        <li>User Authentication and Authorization</li>
                        <li>Contents Management System</li>
                        <li>Realtime Chat Server</li>
                        <li>Image Resizing, Caching, Serving and Uploading</li>
                        <li>and so on</li>
                    </ul>
                </Answer>

                <Question>Why is it important?</Question>
                <Answer>
                    Because these are what we really do, unlike the typical algorithm test.
                </Answer>
            </Col>
        </Row>
    );
};
