import tagData from '@gcsim/data/src/tags.json';
import {model} from '@gcsim/types';
import {useLocation} from 'wouter';

export default function DBEntryTags({
  tags,
}: {
  tags: model.DBTag[] | undefined | null;
}) {
  // const { t: translate } = useTranslation();
  // const t = (key: string) => translate(key) as string; // idk why this is needed

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [_, setLocation] = useLocation();

  return (
    <div className={'flex flex-row overflow-hidden'}>
      {tags
        ?.filter((tag) => tag !== 1)
        .map((tag) => (
          <div
            className="hover:opacity-50 cursor-pointer bg-slate-500 text-xs font-semibold rounded-full px-2 py-1 mr-2 mt-1 whitespace-nowrap "
            key={tag}
            onClick={() => {
              setLocation(`/tag/${tag}`);
            }}>
            {tagData[tag].display_name}
          </div>
        ))}
    </div>
  );
}
