import { useParams } from 'react-router';
import React, { useEffect, useState } from 'react'
import { Box, Button, FormControl, InputLabel, MenuItem, Select } from '@mui/material';
import BACKEND_HOST from '../const';
import { Base64 } from "js-base64";

const KeyPage = () => {

  const { key } = useParams();

  const [content, setContent] = useState([]);

  const fetchKey = async () => {
    const response = await fetch(BACKEND_HOST + '/api/etcd/keys/' + Base64.encode(key));
    const data = await response.json();
    setContent(data.content);
  }

  const [decodeComponent, setDecodeComponent] = React.useState('');

  const [decodeNamespace, setDecodeNamespace] = React.useState('');

  const handleComponentChange = (event) => {
    setDecodeComponent(event.target.value);
  };

  const handleNamespaceChange = (event) => {
    setDecodeNamespace(event.target.value);
  };

  const [decodeData, setDecodeData] = useState({ data: [] });
  const [decodeIsLoading, setDecodeIsLoading] = useState(false);
  const [err, setErr] = useState('');

  const handleClick = async () => {
    setDecodeIsLoading(true);

    const response = await fetch(BACKEND_HOST + '/api/etcd/keys-decode/' + Base64.encode(key) + '?decodeComponent=' + decodeComponent + '&decodeNamespace=' + decodeNamespace);
    if (!response.ok) {
      setErr(err.message);
      setDecodeIsLoading(false);
      return;
    }

    const result = await response.json();

    setDecodeData(result.content);
    setErr('');
    setDecodeIsLoading(false);
  };

  useEffect(() => {
    fetchKey()
  }, [])

  return (
    <div>
      <h1>Key: {Base64.decode(key)}</h1>
      <h1>Content: </h1>
      <p>{content}</p>
      <Box>
        <h2>Decode as</h2>
        <div>
          <FormControl fullWidth>
            <InputLabel id="decode-component-select-label">DecodeComponent</InputLabel>
            <Select
              labelId="decode-component-select-label"
              id="decode-component-select"
              value={decodeComponent}
              label="DecodeComponent"
              onChange={handleComponentChange}
            >
              <MenuItem value='ApiSix'>ApiSix</MenuItem>
              <MenuItem value='Kubernetes'>Kubernetes</MenuItem>
            </Select>
          </FormControl>
          <FormControl fullWidth>
            <InputLabel id="decode-namespace-select-label">DecodeNamespace</InputLabel>
            <Select
              labelId="decode-namespace-select-label"
              id="decode-namespace-select"
              value={decodeNamespace}
              label="DecodeNamespace"
              onChange={handleNamespaceChange}
            >
              <MenuItem value='Minions'>Minions</MenuItem>
              <MenuItem value='Pods'>Pods</MenuItem>
            </Select>
          </FormControl>
        </div>
      </Box>
      {err && <h2>{err}</h2>}
      <Button onClick={handleClick} variant="contained">Decode</Button>
      {decodeIsLoading && <h2>Loading...</h2>}
      {decodeData && <p>{JSON.stringify(decodeData)}</p>}
    </div>
  );
}

export default KeyPage
