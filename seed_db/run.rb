require 'typhoeus'
require 'pry'
require_relative './connection'
require_relative './pull'
require_relative './comment'
require_relative './review'

connection = Connection.new('cleopatchra')

loop do
  current_page = 1
  # TODO: Static URL for now. Will eventually be a CLI that takes an organization/user and a repo name
  # Will also need a Repo model and a repo_id for Pull
  pulls_request = Typhoeus::Request.new('https://api.github.com/repos/facebook/react/pulls', params: { page: current_page })
  pulls_request.run
  pulls_hash = JSON.parse(pulls_request.response.body)

  pull_urls = pulls_hash.map { |pull_hash| pull_hash['url'] }

  pull_urls.each do |pull_url|
    pull = record_pull(pull_url)
    record_comments(pull)
    record_reviews(pull)
  end

  break if pull_urls.length < 30
end

def record_pull(pull_url)
  # I wish there was a better name for this ^_^
  pull_request = Typhoeus::Request.new(pull_url)
  pull_request.run
  pull_hash = JSON.parse(pulls_request.response.body)
  pull = Pull.new(data_hash: pull_hash)
  pull.record
  pull
end

def record_comments(pull)
  # TODO: Pagination
  comments_request = Typhoeus::Request.new("#{pull.url}/comments")
  comments_request.run
  comments = JSON.parse(comments_request.response.body)

  comments.each do |comment_hash|
    comment = Comment.new(data_hash: comment_hash, pull_id: pull.id)
    comment.record
  end
end

def record_reviews(pull)
  reviews_request = Typhoeus::Request.new("#{pull.url}/reviews")
  reviews_request.run
  reviews = JSON.parse(reviews_request.response.body)
  
  reviews.each do |review_hash|
    review_json = review_hash.to_json
    review = Review.new(data_hash: review_hash, pull_id: pull.id)
    review.record

    review_comments_request = Typhoeus::Request.new("#{pull.url}/reviews")
    review_comments_request.run
    review_comments = JSON.parse(review_comments_request.response.body)

    review_comments.each do |comment_hash|
      comment = Comment.new(data_hash: comment_hash, pull_id: pull.id)
      comment.record
    end
  end
end
