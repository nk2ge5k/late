mod components;
extern crate console_error_panic_hook;

use leptos::*;
use leptos_router::*;

use components::*;

fn main() {
    _ = console_log::init_with_level(log::Level::Debug);
    console_error_panic_hook::set_once();

    mount_to_body(|| view! {<App/ >})
}

#[component]
pub fn App() -> impl IntoView {
    log::debug!("rendering <App />");
    on_cleanup(|| {
        log::debug!("cleaning up <App />");
    });

    return view! {
        <div>
            <Router>
                <Header>
                    <HeaderLink to="/">Home</HeaderLink>
                </Header>
                <Main>
                    <Routes>
                        <Route
                            path="/"
                            view=|| view! { <Home /> }
                        />
                    </Routes>
                </Main>
            </Router>
        </div>
    };
}

#[component]
fn Home() -> impl IntoView {
    log::debug!("rendering <Home />");
    on_cleanup(|| {
        log::debug!("cleaning up <Home />");
    });

    return view! {
        <>
            <h1>Hello!</h1>
        </>
    };
}
