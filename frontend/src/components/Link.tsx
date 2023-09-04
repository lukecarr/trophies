import type { FunctionalComponent, JSX } from "preact";

type Props = {
  newTab?: boolean;
};

const Link: FunctionalComponent<JSX.HTMLAttributes<HTMLAnchorElement> & Props> = ({ newTab, children, ...props }) => {
  return (
    <a
      className="font-semibold underline decoration-dotted hover:decoration-solid"
      {...(newTab && { target: "_blank", rel: "noopener" })}
      {...props}
    >
      {children}
    </a>
  );
};

export default Link;
