import { useState } from 'react';
import { Input } from 'antd';
import ReactResizeDetector from 'react-resize-detector';
import { BinaryFileWrapper } from './BinaryFile.styled';

export function BinaryFile({ url, extension, label, onLabel }) {

  const [width, setWidth] = useState(0);

  return (
    <BinaryFileWrapper>
      <ReactResizeDetector handleWidth={false} handleHeight={true}>
        {({ height }) => {
          if (height !== width) {
            setWidth(height);
          }
          return <div style={{ height: '100%', width: width }} />
        }}
      </ReactResizeDetector>
      <div class="player" style={{ width: width, height: width }}>
        <div class="extension">{ extension }</div>
        <div class="label">
          <Input bordered={false} size="small" defaultValue={label} onChange={(e) => onLabel(e.target.value)} />
        </div>
      </div>
    </BinaryFileWrapper>
  )
}

