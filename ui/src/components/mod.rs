use leptos::*;
use leptos_router::*;

#[component]
pub fn Main(children: Children) -> impl IntoView {
    log::debug!("rendering <Main />");
    return view! {
        <main class="container mx-auto pb-20 md:pb-24">{children()}</main>
    };
}

#[component]
pub fn Header(children: Children) -> impl IntoView {
    log::debug!("rendering <Header />");

    return view! {
        <div class="bg-gray-100 border-b pt-4 mb-6">
            <div class="relative overflow-hidden">
                <nav class={[
                  "navigation",
                  "flex",
                  "items-center",
                  "overflow-auto",
                  "left-1",
                  "relative",
                  "min-h-[3.5rem]",
                  "md:container",
                  "md:mx-auto",
                  "md:px-0",
                  "md:-left-3"
                ].join(" ")}>
                    {children()}
                </nav>
            </div>
        </div>
    };
}

#[component]
pub fn HeaderLink<H>(to: H, #[prop(optional)] exact: bool, children: Children) -> impl IntoView
where
    H: ToHref + 'static,
{
    log::debug!("rendering <HeaderLink />");

    let active = {
        let route_path = use_route().path();
        let link_href = to.to_href()();

        route_path == link_href
    };

    return view! {
        <Link
            to={to}
            exact={exact}
            class={[
                "whitespace-nowrap",
                "py-2",
                "group",
                "relative",
                if active {
                    "text-blue-600"
                }else {
                    "text-gray-800"
                },
            ].join(" ")}
        >
            <div
                class={[
                    "px-3",
                    "py-2",
                    "flex",
                    "items-center",
                    "rounded-md",
                    "group-hover:bg-gray-200",
                    "after:absolute",
                    "after:bottom-0",
                    "after:right-3",
                    "after:left-3",
                    "after:h-0.5",
                    "text-text-primary",
                    "after:visible",
                    if active {
                        "after:bg-blue-600"
                    } else {
                        ""
                    },
                ].join(" ")}
            >
                {children()}
            </div>
        </Link>
    };
}

#[component]
pub fn Link<H>(
    to: H,
    #[prop(optional)] exact: bool,
    #[prop(optional, into)] class: Option<AttributeValue>,
    children: Children,
) -> impl IntoView
where
    H: ToHref + 'static,
{
    log::debug!("rendering <Link />");
    return view! {
        <A
            class={class}
            exact={exact}
            href={to}>{children()}</A>
    };
}
