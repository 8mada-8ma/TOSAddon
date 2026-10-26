-- midiflute v0.4.1
-- midiflute.exe が書いた演奏イベントを読み、パイドパイパーの笛を鳴らす。
-- events.txt の形式:
--   1行目 "# <セッション番号>;"
--   以降 "<連番> <種類> <音番号>;"  種類 P=鳴らす S=止める T=遅延計測
-- status.txt: このアドオンの状態を exe に知らせる "<UNIX時刻> <演奏モード 0/1> <3オクターブ目 0/1>;"
-- helper.txt: exe が動いている証拠(exe が2秒ごとにUNIX時刻を書く)
local addonName = "MIDIFLUTE"
local author = "ANIKI"
_G["ADDONS"] = _G["ADDONS"] or {}
_G["ADDONS"][author] = _G["ADDONS"][author] or {}
_G["ADDONS"][author][addonName] = _G["ADDONS"][author][addonName] or {}
local g = _G["ADDONS"][author][addonName]

g.dir = "../addons/midiflute/"
g.inFile = g.dir .. "events.txt"
g.ackFile = g.dir .. "ack.txt"
g.statusFile = g.dir .. "status.txt"
g.helperFile = g.dir .. "helper.txt"
g.langFile = g.dir .. "lang.txt"   -- midiflute.exe で選んだ言語(ja / en / ko)
g.session = g.session or ""
g.lastSeq = g.lastSeq or 0
g.playing = g.playing or {}
g.octave3 = false
g.inMode = false
g.helperOk = nil
g.lastStatusAt = 0
g.lastHelperCheck = 0
-- 1秒あたりに鳴らす音の上限(人の演奏では超えない数。異常な連打でサーバーへ大量に送らないため)
g.maxNotesPerSec = 40
g.rateSec = 0
g.rateCount = 0

-- エラーだけを記録する
local function log(msg)
	local f = io.open(g.dir .. "log.txt", "a")
	if f ~= nil then
		f:write(os.date("%m/%d %H:%M:%S") .. " " .. tostring(msg) .. "\n")
		f:close()
	end
end

local function chat(msg)
	pcall(CHAT_SYSTEM, "[midiflute] " .. msg)
end

local TEXTS = {
	ja = {
		ready = "MIDIキーボードで演奏できます。",
		noexe_short = "midiflute.exe が起動していません",
		noexe = "midiflute.exe が起動していません。起動すると、MIDIキーボードで演奏できます。",
		connected = "midiflute.exe と接続しました。",
		lost = "midiflute.exe との接続が切れました。",
	},
	en = {
		ready = "You can play with your MIDI keyboard.",
		noexe_short = "midiflute.exe is not running",
		noexe = "midiflute.exe is not running. Start it to play with your MIDI keyboard.",
		connected = "Connected to midiflute.exe.",
		lost = "Lost connection to midiflute.exe.",
	},
	ko = {
		ready = "MIDI 키보드로 연주할 수 있습니다.",
		noexe_short = "midiflute.exe가 실행되지 않았습니다",
		noexe = "midiflute.exe가 실행되지 않았습니다. 실행하면 MIDI 키보드로 연주할 수 있습니다.",
		connected = "midiflute.exe와 연결되었습니다.",
		lost = "midiflute.exe와의 연결이 끊어졌습니다.",
	},
}

local function loadLang()
	local f = io.open(g.langFile, "r")
	if f == nil then return end
	local code = string.match(f:read("*a") or "", "^(%a%a)")
	f:close()
	if code ~= nil and TEXTS[code] ~= nil then g.lang = code end
end

local function L(key)
	local t = TEXTS[g.lang or "ja"] or TEXTS.ja
	return t[key] or TEXTS.ja[key]
end

local function checkOctave3()
	local ok, ret = pcall(function()
		return IS_ENABLED_PIED_PIPER_FLUTING_OCTAVE_3(GetMyPCObject())
	end)
	g.octave3 = (ok and ret == 1)
end

local function writeStatus()
	local f = io.open(g.statusFile, "w")
	if f ~= nil then
		f:write(string.format("%d %d %d;", os.time(), g.inMode and 1 or 0, g.octave3 and 1 or 0))
		f:close()
	end
	g.lastStatusAt = os.time()
end

local function helperAlive()
	local f = io.open(g.helperFile, "r")
	if f == nil then return false end
	local t = tonumber(f:read("*a") or "")
	f:close()
	return t ~= nil and math.abs(os.time() - t) <= 6
end

local function readEvents()
	local f = io.open(g.inFile, "r")
	if f == nil then return nil, nil end
	local content = f:read("*a")
	f:close()
	if content == nil then return nil, nil end
	-- 書き込み途中の不完全な内容は ";" が無いので弾かれる
	local session = string.match(content, "^#%s*(%d+);")
	return session, content
end

local function writeAck()
	local f = io.open(g.ackFile, "w")
	if f ~= nil then
		f:write(g.session .. " " .. tostring(g.lastSeq) .. ";")
		f:close()
	end
end

local function stopAll()
	for id, _ in pairs(g.playing) do
		pcall(Piedpiper.ReqStopFluting, id)
	end
	g.playing = {}
end

-- 演奏モードに入った瞬間、それまでに溜まっていた古いイベントを捨てる
local function sync()
	local session, content = readEvents()
	if session == nil then return end
	g.session = session
	local maxSeq = 0
	for seqStr in string.gmatch(content, "(%d+) %a %d+;") do
		local n = tonumber(seqStr)
		if n > maxSeq then maxSeq = n end
	end
	g.lastSeq = maxSeq
	writeAck()
end

local function process()
	local session, content = readEvents()
	if session == nil then return end
	if session ~= g.session then
		-- 今より古いセッション(古い受信プログラムなど)の書き込みは無視する
		if tonumber(g.session) ~= nil and tonumber(session) < tonumber(g.session) then return end
		-- midiflute.exe が再起動された
		g.session = session
		g.lastSeq = 0
	end
	local maxSeq = g.lastSeq
	for seqStr, kind, idStr in string.gmatch(content, "(%d+) (%a) (%d+);") do
		local seq = tonumber(seqStr)
		local id = tonumber(idStr)
		if seq > g.lastSeq then
			if kind == "P" then
				local now = os.time()
				if now ~= g.rateSec then
					g.rateSec = now
					g.rateCount = 0
				end
				g.rateCount = g.rateCount + 1
				-- 未解放の3オクターブ目(25〜36)と、上限を超えた分は鳴らさない
				if id >= 1 and id <= 36 and (id < 25 or g.octave3) and g.rateCount <= g.maxNotesPerSec then
					pcall(Piedpiper.ReqPlayFluting, id)
					g.playing[id] = true
				end
			elseif kind == "S" then
				if g.playing[id] then
					pcall(Piedpiper.ReqStopFluting, id)
					g.playing[id] = nil
				end
			end
			if seq > maxSeq then maxSeq = seq end
		end
	end
	if maxSeq ~= g.lastSeq then
		g.lastSeq = maxSeq
		writeAck()
	end
end

local function onModeOpen()
	checkOctave3()
	sync()
	g.helperOk = helperAlive()
	g.lastHelperCheck = os.time()
	loadLang()
	if g.helperOk then
		chat(L("ready"))
	else
		pcall(ui.SysMsg, L("noexe_short"))
		chat(L("noexe"))
	end
end

function MIDIFLUTE_ON_INIT(addon, frame)
	local ok, err = pcall(function()
		g.frame = frame
		g.inMode = false
		frame:ShowWindow(1)
		frame:RunUpdateScript("MIDIFLUTE_UPDATE", 0.01)
		checkOctave3()
		writeStatus()
	end)
	if not ok then log("ON_INIT エラー: " .. tostring(err)) end
end

function MIDIFLUTE_UPDATE(frame, totalElapsedTime)
	local ok, err = pcall(MIDIFLUTE_UPDATE_BODY)
	if not ok and g.lastErr ~= tostring(err) then
		g.lastErr = tostring(err)
		log("UPDATE エラー: " .. g.lastErr)
	end
	return 1
end

function MIDIFLUTE_UPDATE_BODY()
	local kb = ui.GetFrame("fluting_keyboard")
	local visible = (kb ~= nil and kb:IsVisible() == 1)
	local now = os.time()

	if visible ~= g.inMode then
		g.inMode = visible
		if visible then
			onModeOpen()
		else
			stopAll()
		end
		writeStatus()
	elseif now - g.lastStatusAt >= 1 then
		writeStatus()   -- 1秒ごとに「連携中」を知らせる
	end

	if not visible then return end

	-- 演奏モード中に exe が起動・終了したら知らせる
	if now - g.lastHelperCheck >= 2 then
		g.lastHelperCheck = now
		local alive = helperAlive()
		if alive ~= g.helperOk then
			g.helperOk = alive
			if alive then
				loadLang()
				chat(L("connected"))
			else
				stopAll()
				chat(L("lost"))
			end
		end
	end
	-- midiflute.exe が動いているときだけ演奏データを読む
	if g.helperOk then
		process()
	end
end
