import React, { useState } from 'react'
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormHelperText,
  Grid,
  makeStyles,
} from '@material-ui/core'
import QRCode from 'qrcode.react'
import UserStatusUpdatePreference from '../users/UserStatusUpdatePreference'
import { useSessionInfo } from '../util/RequireConfig'

const useStyles = makeStyles({
  gridItem: {
    width: 'min-content',
  },
})

interface SettingsDialogProps {
  open: boolean
  onClose: () => void
}

// source: w3schools.com/js/js_cookies.asp
function getCookie(cname: string): string {
  let name = cname + '='
  let decodedCookie = decodeURIComponent(document.cookie)
  console.log(document.cookie)
  let ca = decodedCookie.split(';')
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i]
    while (c.charAt(0) == ' ') {
      c = c.substring(1)
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length)
    }
  }
  return ''
}

export default function SettingsDialog(props: SettingsDialogProps) {
  const classes = useStyles()
  const { userID: currentUserID } = useSessionInfo()
  const [showQRCode, setShowQRCode] = useState(false)

  console.log(getCookie('goalert_session.2'))

  return (
    <Dialog
      open={props.open}
      onClose={props.onClose}
      onExited={() => setShowQRCode(false)}
    >
      <DialogTitle>Settings</DialogTitle>
      <DialogContent>
        <Grid container spacing={2}>
          <Grid className={classes.gridItem} item xs={12}>
            {currentUserID && (
              <UserStatusUpdatePreference userID={currentUserID} />
            )}
          </Grid>
          <Grid item xs={12}>
            <Button
              variant='contained'
              color='primary'
              onClick={() => setShowQRCode(true)}
            >
              Link to Mobile App
            </Button>
            <FormHelperText>
              Generates a QR code to be scanned from the GoAlert mobile app to
              log in
            </FormHelperText>
          </Grid>
          {showQRCode && (
            <Grid item xs={12}>
              <QRCode value={getCookie('goalert_session.2')} />
            </Grid>
          )}
        </Grid>
      </DialogContent>
      <DialogActions>
        <Button variant='contained' color='primary' onClick={props.onClose}>
          Done
        </Button>
      </DialogActions>
    </Dialog>
  )
}
