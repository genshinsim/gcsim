import { Editor, type EditorProps } from "@gcsim/components";
import { NonIdealState } from "@gcsim/primitives";

export default (props: EditorProps & { loading: boolean }) => {
	if (props.loading) {
		return <NonIdealState loading />;
	}
	return (
		<div className="w-full 2xl:mx-auto 2xl:container px-2">
			<Editor {...props} showHelpers={false} />
		</div>
	);
};
