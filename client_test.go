//go:generate mockgen -destination=mockprocess_test.go -package=pinentry_test . Process

package pinentry_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/golang/mock/gomock"

	"github.com/sovamorco/go-pinentry/v3"
)

func TestClientClose(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientArgs(t *testing.T) {
	for i, tc := range []struct {
		clientOptions []pinentry.ClientOption
		expectedArgs  []string
	}{
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithArgs([]string{
					"--arg1",
					"--arg2",
				}),
			},
			expectedArgs: []string{
				"--arg1",
				"--arg2",
			},
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithDebug(),
			},
			expectedArgs: []string{
				"--debug",
			},
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithNoGlobalGrab(),
			},
			expectedArgs: []string{
				"--no-global-grab",
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := newMockProcess(t)

			p.expectStart("pinentry", tc.expectedArgs)
			clientOptions := []pinentry.ClientOption{pinentry.WithProcess(p)}
			clientOptions = append(clientOptions, tc.clientOptions...)
			c, err := pinentry.NewClient(clientOptions...)
			assert.NoError(t, err)

			p.expectClose()
			assert.NoError(t, c.Close())
		})
	}
}

func TestClientBinaryName(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry-test", nil)
	c, err := pinentry.NewClient(
		pinentry.WithBinaryName("pinentry-test"),
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientCommands(t *testing.T) {
	for i, tc := range []struct {
		clientOptions   []pinentry.ClientOption
		expectedCommand string
	}{
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithCancel("cancel"),
			},
			expectedCommand: "SETCANCEL cancel",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithDesc("desc"),
			},
			expectedCommand: "SETDESC desc",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithError("error"),
			},
			expectedCommand: "SETERROR error",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithKeyInfo("keyinfo"),
			},
			expectedCommand: "SETKEYINFO keyinfo",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithNotOK("notok"),
			},
			expectedCommand: "SETNOTOK notok",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithOK("ok"),
			},
			expectedCommand: "SETOK ok",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithOption("option"),
			},
			expectedCommand: "OPTION option",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithOptions([]string{
					"option",
				}),
			},
			expectedCommand: "OPTION option",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithPrompt("prompt"),
			},
			expectedCommand: "SETPROMPT prompt",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithQualityBarToolTip("qualitybartooltip"),
			},
			expectedCommand: "SETQUALITYBAR_TT qualitybartooltip",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithTimeout(time.Second),
			},
			expectedCommand: "SETTIMEOUT 1",
		},
		{
			clientOptions: []pinentry.ClientOption{
				pinentry.WithTitle("title"),
			},
			expectedCommand: "SETTITLE title",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := newMockProcess(t)

			p.expectStart("pinentry", nil)
			p.expectWritelnOK(tc.expectedCommand)
			clientOptions := []pinentry.ClientOption{pinentry.WithProcess(p)}
			clientOptions = append(clientOptions, tc.clientOptions...)
			c, err := pinentry.NewClient(clientOptions...)
			assert.NoError(t, err)

			p.expectClose()
			assert.NoError(t, c.Close())
		})
	}
}

func TestClientGetPIN(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	expectedPIN := "abc"
	expectedFromCache := false
	p.expectWriteln("GETPIN")
	p.expectReadLine("D " + expectedPIN)
	p.expectReadLine("OK")
	actualPIN, actualFromCache, err := c.GetPIN()
	assert.NoError(t, err)
	assert.Equal(t, expectedPIN, actualPIN)
	assert.Equal(t, expectedFromCache, actualFromCache)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientConfirm(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	expectedConfirm := true
	p.expectWriteln("CONFIRM confirm")
	p.expectReadLine("OK")
	actualConfirm, err := c.Confirm("confirm")
	assert.NoError(t, err)
	assert.Equal(t, expectedConfirm, actualConfirm)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientConfirmCancel(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	p.expectWriteln("CONFIRM confirm")
	p.expectReadLine("ERR 83886179 Operation cancelled <Pinentry>")
	actualConfirm, err := c.Confirm("confirm")
	assert.Error(t, err)
	assert.True(t, pinentry.IsCancelled(err))
	assert.Equal(t, false, actualConfirm)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientGetPINCancel(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	p.expectWriteln("GETPIN")
	p.expectReadLine("ERR 83886179 Operation cancelled <Pinentry>")
	actualPIN, actualFromCache, err := c.GetPIN()
	assert.Error(t, err)
	assert.True(t, pinentry.IsCancelled(err))
	assert.Equal(t, "", actualPIN)
	assert.Equal(t, false, actualFromCache)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientGetPINFromCache(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	expectedPIN := "abc"
	expectedFromCache := true
	p.expectWriteln("GETPIN")
	p.expectReadLine("S PASSWORD_FROM_CACHE")
	p.expectReadLine("D " + expectedPIN)
	p.expectReadLine("OK")
	actualPIN, actualFromCache, err := c.GetPIN()
	assert.NoError(t, err)
	assert.Equal(t, expectedPIN, actualPIN)
	assert.Equal(t, expectedFromCache, actualFromCache)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientGetPINQualityBar(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	p.expectWritelnOK("SETQUALITYBAR")
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
		pinentry.WithQualityBar(func(pin string) (int, bool) {
			return 10 * len(pin), true
		}),
	)
	assert.NoError(t, err)

	expectedPIN := "abc"
	expectedFromCache := false
	p.expectWriteln("GETPIN")
	p.expectReadLine("INQUIRE QUALITY a")
	p.expectWriteln("D 10")
	p.expectWriteln("END")
	p.expectReadLine("INQUIRE QUALITY ab")
	p.expectWriteln("D 20")
	p.expectWriteln("END")
	p.expectReadLine("INQUIRE QUALITY abc")
	p.expectWriteln("D 30")
	p.expectWriteln("END")
	p.expectReadLine("D abc")
	p.expectReadLine("OK")
	actualPIN, actualFromCache, err := c.GetPIN()
	assert.NoError(t, err)
	assert.Equal(t, expectedPIN, actualPIN)
	assert.Equal(t, expectedFromCache, actualFromCache)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientGetPINQualityBarCancel(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	p.expectWritelnOK("SETQUALITYBAR")
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
		pinentry.WithQualityBar(func(pin string) (int, bool) {
			return 0, false
		}),
	)
	assert.NoError(t, err)

	expectedPIN := "abc"
	expectedFromCache := false
	p.expectWriteln("GETPIN")
	p.expectReadLine("INQUIRE QUALITY a")
	p.expectWriteln("CAN")
	p.expectReadLine("INQUIRE QUALITY ab")
	p.expectWriteln("CAN")
	p.expectReadLine("INQUIRE QUALITY abc")
	p.expectWriteln("CAN")
	p.expectReadLine("D abc")
	p.expectReadLine("OK")
	actualPIN, actualFromCache, err := c.GetPIN()
	assert.NoError(t, err)
	assert.Equal(t, expectedPIN, actualPIN)
	assert.Equal(t, expectedFromCache, actualFromCache)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientGetPINineUnexpectedResponse(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	p.expectWriteln("GETPIN")
	p.expectReadLine("unexpected response")
	actualPIN, actualFromCache, err := c.GetPIN()
	assert.Error(t, err)
	assert.Equal(t, pinentry.UnexpectedResponseError{
		Line: "unexpected response",
	}, err.(pinentry.UnexpectedResponseError)) //nolint:forcetypeassert,errorlint
	assert.Equal(t, "", actualPIN)
	assert.Equal(t, false, actualFromCache)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientReadLineIgnoreBlank(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	p.expectReadLine("")
	p.expectReadLine("\t")
	p.expectReadLine("\n")
	p.expectReadLine(" ")
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func TestClientReadLineIgnoreComment(t *testing.T) {
	p := newMockProcess(t)

	p.expectStart("pinentry", nil)
	p.expectReadLine("#")
	p.expectReadLine("# comment")
	c, err := pinentry.NewClient(
		pinentry.WithProcess(p),
	)
	assert.NoError(t, err)

	p.expectClose()
	assert.NoError(t, c.Close())
}

func newMockProcess(t *testing.T) *MockProcess {
	t.Helper()
	return NewMockProcess(gomock.NewController(t))
}

func (p *MockProcess) expectClose() {
	p.expectWriteln("BYE")
	p.expectReadLine("OK closing connection")
	p.EXPECT().Close().Return(nil)
}

func (p *MockProcess) expectReadLine(line string) {
	p.EXPECT().ReadLine().Return([]byte(line), false, nil)
}

func (p *MockProcess) expectStart(name string, args []string) {
	p.EXPECT().Start(name, args).Return(nil)
	p.expectReadLine("OK Pleased to meet you")
}

func (p *MockProcess) expectWriteln(line string) {
	p.EXPECT().Write([]byte(line+"\n")).Return(len(line)+1, nil)
}

func (p *MockProcess) expectWritelnOK(line string) {
	p.expectWriteln(line)
	p.expectReadLine("OK")
}
