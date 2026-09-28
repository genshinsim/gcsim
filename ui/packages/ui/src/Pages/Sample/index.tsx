import type { Sample } from "@gcsim/types";
import classNames from "classnames";
import Pako from "pako";
import { useEffect, useState } from "react";
import { useDropzone } from "react-dropzone";
import { useTranslation } from "react-i18next";
import SamplePage from "./SamplePage";

export { SamplePage };

export const UploadSample = () => {
	const { t } = useTranslation();
	const [sample, setSample] = useState<Sample | null>(null);
	const [error, setError] = useState<string | null>(null);
	const { acceptedFiles, getRootProps, getInputProps } = useDropzone({
		maxFiles: 1,
		noClick: sample != null || error != null,
	});

	useEffect(() => {
		const file = acceptedFiles[0];
		if (file != null) {
			file.arrayBuffer().then((b64) => {
				try {
					setSample(JSON.parse(Pako.inflate(b64, { to: "string" })));
				} catch (e) {
					let message = "Unknown error when parsing sample...";
					if (e instanceof Error) message = e.message;
					setError(message);
				}
			});
		}
	}, [acceptedFiles]);

	const dzClass = classNames(
		"border-dashed border-2 w-full p-8 flex place-content-center items-center cursor-pointer",
	);

	if (sample == null && error == null) {
		return (
			<div className="p-8">
				<div {...getRootProps({ className: dzClass })}>
					<input {...getInputProps()} />
					<span className="text-g-lg">{t("sample.drop_text")}</span>
				</div>
			</div>
		);
	}

	return (
		<div {...getRootProps({ className: "dropzone" })}>
			<input {...getInputProps()} />
			<SamplePage sample={sample} error={error} />
		</div>
	);
};
