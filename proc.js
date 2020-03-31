const { expose } = require("threads/worker");
const Nextcloudbot = require("nextcloud-talk-bot");

expose(function join() {
  const bot = new Nextcloudbot({
    autoJoin: true
  });

  bot.onText(/^\#(v|V)ideo(call|chat)$/, msg =>
    msg.reply(
      `@${msg.actorId} started a video call. Tap on ${process.env.NCTB_JITSI_URL}/${msg.token} to join!`,
      false
    )
  );

  bot.startPolling();
});
