require 'typhoeus'

class Request
  attr_reader :request, :limit_remaining, :limit_reset_time
  attr_writer :limit_remaining, :limit_reset_time

  def initialize(*args, &block)
    @request = Typhoeus::Request.new(*args, &block)
    @limit_remaining = nil
    @limit_reset_time = nil
  end

  def run
    if limit_remaining === 0
      wait_until_reset
    end
    request.run
    limit_remaining = request.response.headers['X-RateLimit-Remaining'].to_i
    limit_reset_time = request.response.headers['X-RateLimit-Reset'].to_i
  end

  def response
    request.response
  end

  private

  def wait_until_reset
    sleep Time.at(limit_reset_time) - Time.now
  end
end