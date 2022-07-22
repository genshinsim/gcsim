import EnkaToGOOD from "./EnkaToGOOD";

export default async function FetchandValidateDataFromEnka(validUid: string) {
  const enkaResponse = await fetch(`https://enka.gcsim.app/${validUid}`);
  if (!enkaResponse.ok) {
    throw new Error(`Failed to fetch ${validUid}`);
  }

  const enkaData = await enkaResponse.json();
  return EnkaToGOOD(enkaData);
}
