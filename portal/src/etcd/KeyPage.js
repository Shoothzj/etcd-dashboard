import { useParams } from 'react-router';
import React, { useEffect, useState } from 'react';
import { Box, Button, FormControl, InputLabel, MenuItem, Select } from '@mui/material';
import { Base64 } from 'js-base64';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogActions from '@mui/material/DialogActions';
import DialogTitle from '@mui/material/DialogTitle';
import BACKEND_HOST from '../Const';

function KeyPage() {
  const { key } = useParams();

  const [content, setContent] = useState([]);

  const [hexContent, setHexContent] = useState([]);

  const [open, setOpen] = React.useState(false);
  const fetchKey = async () => {
    const response = await fetch(`${BACKEND_HOST}/api/etcd/keys/${key}`);
    const data = await response.json();
    setContent(data.content);
  };

  const fetchHexKey = async () => {
    const response = await fetch(`${BACKEND_HOST}/api/etcd/keys/${key}?codec=hex`);
    const data = await response.json();
    setHexContent(data.content);
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

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

    const response = await fetch(
      `${BACKEND_HOST}/api/etcd/keys-decode/${key}?decodeComponent=${decodeComponent}&decodeNamespace=${decodeNamespace}`,
    );
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
    fetchKey();
    fetchHexKey();
  }, []);

  return (
    <div>
      <h1>
        Key:
        {Base64.decode(key)}
      </h1>
      <Button variant="contained" onClick={handleClickOpen}>
        show hex
      </Button>
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description">
        <DialogTitle id="alert-dialog-title">hex content</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">{hexContent}</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} autoFocus>
            close
          </Button>
        </DialogActions>
      </Dialog>
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
              onChange={handleComponentChange}>
              <MenuItem value="ApiSix">ApiSix</MenuItem>
              <MenuItem value="Kubernetes">Kubernetes</MenuItem>
            </Select>
          </FormControl>
          <FormControl fullWidth>
            <InputLabel id="decode-namespace-select-label">DecodeNamespace</InputLabel>
            <Select
              labelId="decode-namespace-select-label"
              id="decode-namespace-select"
              value={decodeNamespace}
              label="DecodeNamespace"
              onChange={handleNamespaceChange}>
              <MenuItem value="Minions">Minions</MenuItem>
              <MenuItem value="Pods">Pods</MenuItem>
            </Select>
          </FormControl>
        </div>
      </Box>
      {err && <h2>{err}</h2>}
      <Button onClick={handleClick} variant="contained">
        Decode
      </Button>
      {decodeIsLoading && <h2>Loading...</h2>}
      {decodeData && <p style={{ whiteSpace: 'pre-wrap' }}>{decodeData.toString()}</p>}
    </div>
  );
}

export default KeyPage;
