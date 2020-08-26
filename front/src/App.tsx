import React                                from "react";
import "antd/dist/antd.dark.min.css";
import {Layout, Menu}                       from "antd";
import {MainPage}                           from "./page/MainPage";
import {BrowserRouter, Link, Route, Switch} from "react-router-dom";
import {ProblemListPage}                    from "./page/ProblemListPage";
import {ProblemPage}                        from "./page/ProblemPage";

const {Header, Footer, Content} = Layout;

const RouteTable: { path: string, component: any, menuName?: string }[] = [
    {path: "/", component: <MainPage/>, menuName: "🍄 Bonitto"},
    {path: "/problems", component: <ProblemListPage/>, menuName: "Problems"},
    {path: "/problems/:no", component: <ProblemPage/>},
];

function App() {
    return (
        <BrowserRouter>
            <Layout>
                <Header>
                    {/*TODO : pathname reactively*/}
                    <Menu theme="dark" mode="horizontal" defaultSelectedKeys={[window.location.pathname]}>
                        {RouteTable.filter((r) => r.menuName).map((r) =>
                            <Menu.Item key={r.path}>
                                <Link to={r.path}>{r.menuName}</Link>
                            </Menu.Item>,
                        )}
                    </Menu>
                </Header>
                <Content style={{padding: 60}}>
                    <Switch>
                        {RouteTable.map((r) =>
                            <Route exact path={r.path} key={r.path}>{r.component}</Route>,
                        )}
                    </Switch>
                </Content>
                <Footer style={{textAlign: "center"}}>
                    <p>Bonitto &copy;2020 | <a href="https://eunchan.com">Eunchan Lee</a></p>
                </Footer>
            </Layout>
        </BrowserRouter>
    );
}

export default App;
